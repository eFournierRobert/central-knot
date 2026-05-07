# Central Knot

Central knot is an HTTP BitTorrent tracker that follows the specs described in [BEP 3](https://bittorrent.org/beps/bep_0003.html) and 
[BEP 23](https://bittorrent.org/beps/bep_0023.html). It has some extra events to make it work
with qBittorrent, since it is the client I tested it with.

## How it works
When launched, it creates a Sqlite database and serves the endpoint `/announce` over the port 9000.
After that, adding `http://<central-knot address>:9000/announce` to a torrent trackers will put it
in use.

The BEP doesn't really specify error handling other than a response with the key `failure reason`,
so it is how Central knots handles them.