

Add SQL to store table of network topos, and the generated secret that the HTTPS API will require
THATS IT No more data. no users. no badger no vouchers yet

HTTPS SUBMISSIONS:
    /event

HTTPS  API:
    /api/cnc  ------- All UI requires POST here a JSON encoded command
    /api/action ----- All action scripts can send commands/requests here
    /api/env      --- K/V Store for authd keys



UI will be a seperate running program that uses the databse file to load the secret required to
edit the network and apply changes.

UI will NEVER write to the database directly. It simply reads the information required to control
the server.





