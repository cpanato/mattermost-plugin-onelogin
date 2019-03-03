# OneLogin Notification [![CircleCI](https://circleci.com/gh/cpanato/mattermost-plugin-onelogin.svg?style=svg)](https://circleci.com/gh/cpanato/mattermost-plugin-onelogin)

This plugin receive and post notifications from OneLogin Webhook

For now it only parses the user login and check the configured risk threshold and then post a message in the specified channel

## Configuration

#### Mattermost side

    - Install the plugin
    - Configure the plugin:
        - in the `TeamChannel` field add a the team and the channel you want to post the messages separated by comma. ie. `TeamA,ChannelX`.
        - set the `RiskThreshold` the value is from 0 to 100.
        - set the `Username` which is the user the will be use to post the messages.
        - set the `Token` this will be used to set the webhook header in the OneLogin side in order to validate the request.

#### OneLogin side

    - Create the webhook and set the Header `X-OneLogin-Token` with the value you created in the Mattermost configuration, see above.



## Next features

 - add slash command to block user
 - parse other events types