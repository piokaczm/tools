# Spreadlove

It's an application for spreading love via bonus.ly - as previous idea of doing lottery using bonus.ly coins is a bit sketchy I had to come up with something equally fun but less suspicious.

It basically sends a bonus to some folks on bonus.ly, but you can set flags to make it more precise (or random...). So using different flags combinations you can split your remaining coins/specific amount of coins between specific colleagues/random colleagues.

Code is not very clean and it lacks tests as it was coded during lunch brake, so feel free to contribute and make it better.

#### Prerequisites

To use it you have to get your own bonus.ly API key and have it as `BONUSLY_API_KEY` env variable.

#### How to use it?

After compiling it by yourself or adding one of provided binaries to your `PATH` just run:
```
spreadlove <flags>
```

where flags are as following:

`lucky-folks=<int>` - to how many random folks do you want to send some shiny coins?
`lucky-names=<strings separated by a comma without spaces>` - if you want to give TCs to a specific group of people instead of random folks, just provide this flag.
`coins-limit=<int>` - if you don't want to spend ALL your coins at once you can specify a limit using this one.
`message=<string>` - if you want to add your custom message for a bonus, provide it here.

You can also use `./install` to install the tool to your $GOBIN and copy autocompletion script to osx's default bash_completion.d directory.

#### Autocompletion

There's a script for autocompletion living under `./bin/autocompletion`, add it to your `bash_completion.d/spreadlove` and have fun spreading your coins even faster.