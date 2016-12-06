#!/usr/bin/perl -w
#
use LWP::Simple;

my $server = 'www.dgears.com:1337';
my $uri = '/blacklistme/api';
my $urlbase = 'http://'.$server.$uri;
my $apikey = '8078ac0a68877d828efccb68e91dabc1e720a716acb92d71e4f278428ec2a311';
my $email = 'alane@fueledcafe.com';

my $url = $urlbase.'?email='.$email.'&apikey='.$apikey;
my $response = get $url;
die 'Error getting $url' unless defined $response;
print $response."\n";

