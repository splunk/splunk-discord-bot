# Splunk Discord Bot

This bot can be added to any Discord server, and will record messages by sending them to Splunk.

# Setup

## Create your own application and bot
Follow the first directions here to create your app and register your bot:

https://golangexample.com/discord-bot-in-golang/

Set up your bot and add it to your server.

## Set up config.json

Set up in config.json the following:

```json lines
{
  "token": "your bot token", 
  "hec_endpoint": "the splunk host",
  "hec_token": "the HEC token",
  "hec_index": "the index to send to",
  "insecure_skip_verify": false, # whether to allow insecure TLS connections to Splunk. 
}
```

## Run it locally

```bash
$> make run
```

# License

Copyright 2022 Splunk, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
