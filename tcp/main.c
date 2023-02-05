#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>

int main(int argc, char const *argv[])
{
    signal(SIGPIPE, SIG_IGN);

    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0)
    {
        perror("socket error");
        exit(EXIT_FAILURE);
    }

    int enable = 1;
    int res = setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &enable, sizeof(int));
    if (res < 0)
    {
        perror("setsockopt error");
        exit(EXIT_FAILURE);
    }

    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(8888);

    socklen_t addrlen = sizeof(addr);
    res = bind(fd, (struct sockaddr *)&addr, addrlen);
    if (res < 0)
    {
        perror("bind error");
        exit(EXIT_FAILURE);
    }

    res = listen(fd, 3);
    if (res < 0)
    {
        perror("listen error");
        exit(EXIT_FAILURE);
    }

    char *ok = "ok\n";
    while (1)
    {
        printf("waiting for connection...\n");
        int accepted = accept(fd, (struct sockaddr *)&addr, (socklen_t *)&addrlen);
        if (accepted < 0)
        {
            perror("accept error");
            exit(EXIT_FAILURE);
        }
        printf("...accepted connection: %d\n", accepted);

        while (1)
        {

            printf("reading...\n");
            char buffer[1024] = {0};
            int n = read(accepted, buffer, 1024);
            printf("...read: %d:\n", n);
            if (n <= 0)
            {
                break;
            }

            printf("sending...\n");
            n = send(accepted, ok, strlen(ok), 0);
            printf("...sent: %d\n", n);
            if (n <= 0)
            {
                break;
            }
        }

        printf("closing...\n");
        close(accepted);
    }

    shutdown(fd, SHUT_RDWR);
    return 0;
}