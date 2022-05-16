# Release and Upgrade procedures and communications

The `#validators-private` channel on discord will be used for all communications from the team. Only **active validators** should be allowed access, for security reasons.

**The core team will endeavour to always make sure there is 48-72 hours notice of an impending upgrade, unless there is no alternative.**

Most of our validator communications is done on the [Terra Validator Discord](https://discord.com/invite/xfZK6RMFFx). You should join, and change your server name to `nick | validator-name`, then ask a mod for permission to see the private validator channels.

## Release versioning

**If a change crosses a major version number, i.e. `1.x.x -> 2.x.x` then it is definitely consensus-breaking.**

In the past, some releases have been consensus-breaking but only incremented a minor version, if clearly indicated. In future we will look to be clearer. 

**Only patch versions, i.e. `x.x.1 -> x.x.2`, or `1.1.0 -> 1.1.1` are guaranteed to be non-consensus breaking.**

## Scheduled upgrade via governance

For a SoftwareUpgradeProposal via governance:

1. Validators will be told via announcements channel when the prop is live
2. Validators will be told via announcements channel if it passes
3. Validators will be told via announcements channel when the upgrade instructions are available, and the upgrade will be coordinated in the private validators channel as the target upgrade block nears.

## Emergency upgrade or security patch

If the team needs to upgrade the chain faster than the cadence of governance, then a special procedure applies.

This procedure minimizes the amount that is publicly shared about a potential issue.

1. An announcement calling validators to check in on the private validators channel will be posted on the _validators announcement channel_ on discord. No specifics will be shared here, as it is public.
2. Details of the patch and the upgrade plan will be shared on the private channel, as well as an expected ETA.
3. When instructions are available, they will be pinned, and a second announcement sent on the announcements channel. A thread for acknowledgements will be created for validators to signal readiness.
4. The team will compile a spreadsheet of validator readiness to check we are past 67%.

There are two further considerations:

1. If the change is consensus-breaking, a halt-height will be applied and validators will manually upgrade at that block, after halt.
2. If the change is non-consensus breaking, validators will apply when ready, and then signal readiness.

## Syncing from genesis

The team will be putting together instructions that will be kept up-to-date for syncing without using a backup.
