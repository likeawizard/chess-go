# Chess engine written in Golang
## Installation
* Clone repository or copy files
* Make sure you have Golang 1.17+ installed
* On linux systems simply run `./build.sh` or manually compile each of the endpoints under `./cmd/**`
* You should have three separate executables `analyze, self-play, lichess`
* `cp .env.dev .env`

## Setting up .env
* For `./self-play` you can set `STARTING_FEN` to set a custom starting position. (Comment out to use standard starting position)
* Set search depth `EVALUATION_DEPTH` number of plies (half-moves) to look into. For example `EVALUATION_DEPTH=2` will look 2 plies deep- next move for the side to play, response from opposition and another move for current side to move.
* Set `EVALUATION_ALGO=alphabete(default)|minmax` to set what search algorithm will be used.
* For `./lichess` only. Set `LICHESS_TOKEN` to your bot account token. **Do not publish this token as it is equivalent to your password.** To learn more about how to set up a lichess bot account see: [Upgrading To a Lichess Bot Account](https://lichess.org/api#operation/botAccountUpgrade)


## Usage
### Self Play
* Set up your `.env` file
* Run `./self-play`
* Sending an interrupt signal `Ctrl+c` will also output the full movelist of the game before shutdown
### Analyze
* Run `./analyze -fen="<position to analyze>"`
### Lichess
* Run `./lichess`.
* The bot will authenticate using the token from `.env` and listen for incomming challenges and moves