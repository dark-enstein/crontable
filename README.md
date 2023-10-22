# crontable
a cron expression parse

#### i guess this is my attempt to make a parser within 20 hours, didn't quite work out completely well, right now the `reader` and `meaning` pkgs are a glass house, sitting carefully behing their test cases, any variation comes like a small pebble just heavy enough. will try sometime else; hopefully I'm more experienced in the art; but anyways here below is how it should be used

```
-> crontable <path to crontab>
Your crontab expression reads,
"On the 5th and 9th minute, between the 3rd and 19th hour, every 5th of the month, on the 5th month, every day of the week"
```