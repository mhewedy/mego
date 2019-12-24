# MEGO

## What?
A meeting organizer for Exchange server.

## Why?
Reserving a meeting with many people and many candidate rooms wasn't always an amusing activity. you need go by eye through the timeline of all people and all candidate rooms to see the most suitable one.

In MEGO, all people with each selected room is **represented in just one box**, either it is *free* or *not*. (busy/tentative)

so you no longer need to do the dunting work of eye checking, and adding and removing rooms by hand, all rooms should be added in a tree-like structure, organized around the building, the zone, the size, so you just click the room(s) and choose the users and here we go.

Other features also exists like easily search through your organization people through a [simple and efficient indexing algorithm](https://medium.com/@mhewedy_46874/implementing-a-simple-indexing-algorithm-in-golang-c65be7eaa563)

## How it works:
* It connects to Exchage server using [ews](http://github.com/mhewedy/ews)
* The configuration is defined in [app.conf](https://github.com/mhewedy/mego/blob/master/mego-api/app.conf), so you need to place it in the same directory as the binary and configure it properly
* Room list is defined in [rooms.csv](https://github.com/mhewedy/mego/blob/master/mego-api/rooms.csv), so you need to place it in the same directory as the binary and fill the values according to your environment.
* The user credentials entered at login screen is encryped and saved in memory with an encryption key itself is defined only when the server starts. (see [enc.go](https://github.com/mhewedy/mego/blob/master/mego-api/user/enc.go))    
So when a request made to the server with the valid user token, then the token got verified first then the `username` is being exctracted and used to get the encryped password from the in-memory user map and then the password got decrypted by the key defined at the server startup.

## Screenshots

1. Search is done using some basic indexing and searching in-memory scheme.
<kbd> <img src="https://github.com/mhewedy/mego/raw/master/screenshots/1.png"> </kbd>

2. Availabity is returned based on the user input
<kbd> <img src="https://github.com/mhewedy/mego/raw/master/screenshots/2.png"></kbd>

3. User can send meeting request and add optional attendees as well, where thier availabity time will not checked.
<kbd> <img src="https://github.com/mhewedy/mego/raw/master/screenshots/3.png"></kbd>

## Download
Download the latest release from [the releases section](https://github.com/mhewedy/mego/releases/latest) (Linux, mac and win releases available) 

