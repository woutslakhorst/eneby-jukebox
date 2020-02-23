# Eneby jukebox

My 4-year old asked if he could play some music in his room....

That used to be the sign to buy a my-first-sony, walkman, discman, mp3-player or iPod. But these days everything is streaming or involves a phone or tablet.
Giving a phone or tablet to my children and configuring it to be child-safe sounded like a lot of work. Plus the music that comes from it is kinda crap.

So, there must be something we can hack together (quickly). We already had some playlists in place on Spotify, so playing those would be kinda nice. After some thinking and investigating: raspberry-pi (1st goto when you need some online device) combined with some rfid cards and a reader and a Eneby speaker.

The Eneby speaker actually came later in the plan, the first plan was to buy some cheap speakers and hack them into some container. But the sound is awful.
Then I came across the Eneby speaker from Ikea, costing just euro 50 and delivering quite a nice sound (for its price).

## Contents

- rfid-mopidy-go: Go code for listening to the rfid reader and talking to the Mopidy API.
- 50-usb-rfid-reader.rules: rules file for raspberry in /etc/udev/rules.d/

## Status led

Enable serial port by editing /boot/config.txt

```
enable_uart=1
```

Connect pin 8 (TXD) and pin 9 (Ground) (https://pinout.xyz/) to a led.

I used a green led, therefore 180Ω has to be used. (calculate yours with: https://www.ledtuning.nl/nl/resistor-calculator)  

## On/Off button

shorting pin 5 and 6 will power-up the pi from hibernation. The trick to to put the pi in hibernation (halt command) by using the same pins

It's now built into the rfid-mopidy-go executable using https://github.com/stianeikeland/go-rpio (had to fork a PR to get PI4 fix)

apparently:

adding to /boot/config.txt
```
dtoverlay=gpio-shutdown,gpio_pin=3
```

and disabling i2c via rasp-config

## Some UX

The status led comes on quite early, playing a sound when services have started could help...