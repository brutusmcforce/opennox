#ifndef NOX_MATH_H
#define NOX_MATH_H

// Math functions for the game logic. The C library implementations differ between platforms in the last bits
// (glibc and msvcrt return different results for about 10% of sin/cos inputs), which makes the simulation
// diverge between the Linux and Windows builds. These use Go's math package, which is the same everywhere.
double nox_math_sin(double x);
double nox_math_cos(double x);
double nox_math_acos(double x);
double nox_math_atan2(double y, double x);
double nox_math_pow(double x, double y);

#endif // NOX_MATH_H
