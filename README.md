# Sandy

> A tiny sandbox to run untrusted code. 🏖️

Sandy uses Ptrace to hook into READ syscalls, giving you the option to accept or deny syscalls before they are executed.

## Usage

```
Usage of ./sandy:

  sandy [FLAGS] command

  flags:
    -h	Print Usage.
    -n value
        A glob pattern for automatically blocking file reads.
    -y value
        A glob pattern for automatically allowing file reads.
```

## Use cases

### You want to install anything

```shell
> sandy -n "/etc/password.txt" npm install sketchy-module

  BLOCKED READ on /etc/password.txt
```

```shell
> sandy -n "/etc/password.txt" bash <(curl  https://danger.zone/install.sh)

  BLOCKED READ on /etc/password.txt
```

### You are interested in what file reads you favourite program makes.

Sure you could use strace, but it references file descriptors sandy makes the this much easier at a glance by printing the absolute path of the fd.

```
> sandy ls
Wanting to READ /usr/lib/x86_64-linux-gnu/libselinux.so.1 [y/n]
```

### You _don't_ want to buy your friends beer

A friend at work knows that you are security conscious and that you keep a `/free-beer.bounty` file in home directory. With the promise of a round of drinks and office wide humiliation Dave tries to trick you with a malicious script under the guise of being a helpful colleague.

You run there script with sandy and catch him red handed.

```shell
> sandy -n *.bounty bash ./dickhead-daves-script.sh

  BLOCKED READ on /free-beer.bounty
```

**NOTE**: It's definitely a better idea to encrypt all your sensitive data, sandy should probably only be used when that is inconvenient or impractical.

**NOTE**: I haven't made any effort for cross-x compatibility so it currently only works on linux. I'd happily accept patches to improve portability.
