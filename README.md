# avr-demo
Demo programs using edmccard/avr-sim

## avrbeep
Plays a beep (requires PortAudio).

Simulates an ATmega8 @ 1 Mhz with a speaker connected to pin 0 of PORTB,
running BEEP.ASM from https://sites.google.com/site/avrasmintro/

## nucleik
Plays a song (requires PortAudio).

Simulates an ATmega8 @ 1Mhz with a speaker connected to pin 0 of PORTB.
The code is a custom port to AVR assembly from the 6502 assembly of the
"Nucleik" Apple II demo
(disk 2 of [maxagaz.shk](http://mirrors.apple2.org.za/ground.icaen.uiowa.edu/apple8/Music/maxagaz.shk))

## avrstdio
"Terminal" input and output.

Simulates an ATmega8 communicating via USART with an RS-232 terminal.

