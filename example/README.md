# Example

This example sets up a local Discord bot and a Splunk Enterprise instance.

The example has two different use cases: sending chats to Splunk, and sending alerts to Discord.

## Prerequisites
You will need to install Docker Compose 1.29 or higher.

## Set up
Add a token to the discord bot:
* copy config.json.example into config.json
* set the token field value to the bot token (see instructions in the root folder of the repository)

Set up the webhook:
* Create a private channel where you will add the bot and copy its ID
* Paste the channel ID under the foo webhook

Start the docker compose deployment:
`$>  docker-compose up --build`

## Send discord chats to Splunk
Add your bot to Discord and add it to a channel you want to listen to.

Send a few messages and check the main index of Splunk.
* Authenticate to http://localhost:18000 with `admin`/`changeme`.
* Go to the search application, and enter `index=main`

## Send Splunk alerts to Discord
You can set up Splunk to send alerts to Discord over a webhook.

First, create an alert in Splunk:
* Authenticate to http://localhost:18000 with `admin`/`changeme`.
* Go to the search application, and enter `index=main boo`
* Save this search as an alert. Make it a real time alert.
* In the alert action, pick webhook and enter the following: `http://bot:8080/?webhook=foo`

Trigger the alert.
* Type boo in the Discord channel

Check Discord to see the alert rendered.
* You will see a message such as:
```boo - see results http://27014f39e856:8000/app/search/search?q=%7Cloadjob%20rt_scheduler__admin__search__boo_at_1655880844_6.5%20%7C%20head%201%20%7C%20tail%201&earliest=0&latest=now```
