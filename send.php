<?php
$token = "#";
$chat_id = "#";
$log = "#";

if (!file_exists($log)) {
    die("Log tidak ditemukan: $log");
}

$baca = file($log, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
$tgl = "";
$hasil = [];

foreach ($baca as $isi) {
    if (preg_match('/^\s*(\d{4}-\d{2}-\d{2})/', $isi, $cocok)) {
        $tgl = $cocok[1];
        $hasil = [];
    }
    if (strpos($isi, "Successfully reseted all servers quota") !== false) {
        $hasil[] = "$tgl - $isi"; 
    }
}

if (empty($hasil)) {
    die("Tidak ada log");
}

$pesan = "Log Terbaru:\n\n" . implode("\n", $hasil);

$url = "https://api.telegram.org/bot$token/sendMessage";
$data = [
    "chat_id" => $chat_id,
    "text" => $pesan,
    "parse_mode" => "Markdown"
];

$opt = [
    "http" => [
        "header" => "Content-type: application/x-www-form-urlencoded",
        "method" => "POST",
        "content" => http_build_query($data)
    ]
];

$ctx = stream_context_create($opt);
$res = file_get_contents($url, false, $ctx);

if (!$res) {
    die("Gagal kirim ke Telegram");
}

echo "Berhasil kirim ke Telegram!";
