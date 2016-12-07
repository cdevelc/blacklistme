<?php

$server = 'www.dgears.com:1337';
$uri = '/blacklistme/api';
$urlbase = 'http://'.$server.$uri;
$apikey = '8078ac0a68877d828efccb68e91dabc1e720a716acb92d71e4f278428ec2a311';
$email = 'alane@FUELEDCAFE.com';

$url = $urlbase.'?email='.$email.'&apikey='.$apikey;
$response = file_get_contents($url);
$jresponse = json_decode($response, true);
print_r( $jresponse);
?>