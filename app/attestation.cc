#include <iostream>
#include <string>
#include <stdlib.h>

int main() {
    char bInLen[4];
    std::cin.read(bInLen, 4);
    unsigned int inLen = *reinterpret_cast<unsigned int *>(bInLen);
    char *inMsg = new char[inLen];
    std::cin.read(inMsg, inLen);
    std::string inStr(inMsg);
    delete[] inMsg;

    std::string message = "{\"message\":\"test\"}";
    unsigned int length = message.length();

    std::cout << char(length >> 0)
        << char(length >> 8)
        << char(length >> 16)
        << char(length >> 24);
    std::cout << message << std::flush;
}
