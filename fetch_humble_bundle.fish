#!/usr/bin/env fish

# usage: fetch_humble_bundle.fish ORDERID PATHTOCOOKIE

set orderid $argv[1]
set cookie (cat $argv[2])

# API called at https://www.humblebundle.com/downloads?key=ORDERID
# Grab the _simpleauth_sess cookie from the browser's API call
if ! test -e order_$orderid.json
  curl "https://www.humblebundle.com/api/v1/order/$orderid?wallet_data=true&all_tpkds=true" \
    -H 'DNT: 1' \
    -H "Cookie: _simpleauth_sess=$cookie" \
    > order_$orderid.json
end

for dl in (jq -c '.subproducts[].downloads[].download_struct[]' order_$orderid.json)
  set dir (echo $dl | jq -r '.name')
  set sha1 (echo $dl | jq -r '.sha1')
  set md5 (echo $dl | jq -r '.md5') # Not all have SHA1, just do both
  set web_url (echo $dl | jq -r '.url.web')
  set torrent_url (echo $dl | jq -r '.url.bittorrent')
  # Torrent has no filename via `wget --content-disposition`, so extract from URL instead
  set web_filename (string replace -r '^.+/([^/?]+)\?.+$' '$1' $web_url)
  set torrent_filename (string replace -r '^.+/([^/?]+)\?.+$' '$1' $torrent_url)

  mkdir -p $dir Torrents/$dir
  test -e $dir/$web_filename || wget $web_url -O $dir/$web_filename
  test -e Torrents/$dir/$torrent_filename || wget $torrent_url -O Torrents/$dir/$torrent_filename
  echo "$sha1  $dir/$web_filename" >> sums.sha1
  echo "$md5  $dir/$web_filename" >> sums.md5
end

sha1sum -c sums.sha1
md5sum -c sums.md5
