# C8-compiler

c8-compiler is a compiler that converts code written in c8-lang to an extended version of chip-8 hexadecimal opcodes. The compiler is written in the Go programming language.

## What is Chip-8?

[Chip-8](https://en.wikipedia.org/wiki/CHIP-8) is an interpreted programming language that was designed in the 1970s for the COSMAC VIP and Telmac 1800 computers. It was later used on several other computers and gaming consoles, and is still popular among hobbyist programmers today. Chip-8 has a very simple instruction set, consisting of 35 opcodes that allow basic operations such as drawing graphics, playing sound, and reading input.

## C8-lang

C8-lang is a new programming language that is designed to be compiled to an extended version of Chip-8 opcodes. This version contains two additional opcodes not present in the original set of instructions:

- `9XY1`: save `VX` in the first 8 bits of `I` and `VY` in the last 8 bits.
- `9XY2`: save the first 8 bits of `I` in `VX`, and the last 8 bits in `VY`.

These opcodes were added to increase the functionality of the language and allow for more complex programs to be written.

There are nine primitive functions available in c8-lang:

1. `draw(x, y, length, sprite)`: Receives as parameters three bytes and a pointer to byte. The first byte represents the x coordinate in which the draw is going to be set, the second one represents the y coordinate, the third one represents the length of the sprite, and the fourth represents the address of the sprite (ideally a pointer to the first element of an array of bytes, where the amount of elements represents the height of the sprite, and each bit a pixel in the screen). It returns a bool that is true only when a collision happens.

2. `clean()`: Doesn't receive any parameters or return anything. It cleans the screen.

3. `setDT(value)`: Receives a byte and doesn't return anything. It changes the value of the delay timer.

4. `setST(value)`: Receives a byte and doesn't return anything. It changes the value of the sound timer.

5. `getDT()`: Doesn't receive any parameters and returns the delay timer (a byte).

6. `random()`: Doesn't receive any parameters and returns a random byte.

7. `waitKey()`: Doesn't receive any parameters. Waits until a key is pressed and returns the key value (a byte).

8. `isKeyPressed(key)`: Receives a byte as a parameter and returns a boolean that is true only if the key received as parameter was pressed.

9. `drawFont(x, y, value)`: Receives three bytes as parameters. The first represents the x coordinate of the draw, the second represents the y coordinate, and the third must be a byte between 0 and 15. It draws the character corresponding to that byte at the specified location (x, y).



## Custom Chip-8 Emulator

In order to use the new extended version of Chip-8 opcodes, a custom emulator is needed.
[Here](https://github.com/NoetherianRing/Chip-8) is a custom emulator that supports the additional opcodes.

## Examples

An example of c8-lang code can be found in the `examples` folder of this repository. This code can be compiled using the c8-compiler to produce a Chip-8 program that can be run on a Chip-8 emulator that supports the additional opcodes.

## Grammar

The definition of the c8-lang language can be found in the `grammar.txt` file in this repository. This file contains a formal grammar specification for the language, which can be used to understand the syntax of c8-lang code and to create new programs in the language.

## Usage

To use c8-compiler, you will need to have the Go programming language installed on your system. Once you have Go installed, you can clone this repository and run the `c8-compiler` command to compile your c8-lang code into a Chip-8 program.

The `c8-compiler` command takes two arguments:

1. The first argument is the name of the input file that contains the c8-lang source code you want to compile.
2. The second argument is the name and address of the output ROM file that will be created.

Here is an example of how to use `c8-compiler` to compile the `example.txt` source code in the `examples` folder and create a new ROM file called `example.ch8` in the same folder:
```shell
$ go get github.com/NoetherianRing/c8-compiler
$ git clone https://github.com/NoetherianRing/c8-compiler.git
$ cd c8-compiler
$ go build
$ ./c8-compiler examples/example.txt examples/example.ch8

```

This will produce a file named `example.ch8` that contains the compiled program. You can then run this file on a Chip-8 emulator that supports the extended opcodes to see the program in action.

Note that the ROM files should be used in Chip-8 emulators with more memory than the original one, in order to accommodate the necessities of c8-lang.

## License

c8-compiler is licensed under the MIT License. See the `LICENSE` file in this repository for more information.
