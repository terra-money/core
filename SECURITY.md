# Terra bug reporting and feature requests

The Terra core development team uses GitHub to manage feature requests and bugs. This is done via GitHub Issues.

## Triage and progress 🔜

Issues added to GitHub will be triaged as they come in.

Tracking of in-flight issues will be done through the Terra Core project board, but of course we reserve the right to not make a public issue if there is a security implication in doing so.

## Feature request 🚀

For a feature request, e.g. module inclusion, please make a GitHub issue. Clearly state your use case and what value it will bring to other users or developers on Terra.

If it is something that can be handled by a param change, discuss it on Discord in the `#governance` channel, and consider a governance proposal.

## Standard priority bug 🐛

For a bug that is non-sensitive and/or operational in nature rather than a critical vulnerability, please add it as a GitHub issue.

If it is not triaged in a couple of days, feel free to tag `@the-frey` or `@jakehartnell`.

## Critical bug or security issue 💥

If you're here because you're trying to figure out how to notify us of a security issue, go to [Discord](https://discord.com/invite/xfZK6RMFFx), and alert the core engineers:


- Yun `Yun#4870`
- Paul `Paul Kim#9433`
- Jesse ` jess#5659`

Please avoid opening public issues on GitHub that contain information about a potential security vulnerability as this makes it difficult to reduce the impact and harm of valid security issues.

### Coordinated Vulnerability Disclosure Policy

We ask security researchers to keep vulnerabilities and communications around vulnerability submissions private and confidential until a patch is developed. In addition to this, we ask that you:

- Allow us a reasonable amount of time to correct or address security vulnerabilities.
- Avoid exploiting any vulnerabilities that you discover.
- Demonstrate good faith by not disrupting or degrading Terra’s network, data, or services.

### Vulnerability Disclosure Process

Terra uses the following disclosure process:

- Once a security report is received, the Terra core development team works to verify the issue.
- Patches are prepared for eligible releases in private repositories.
- We notify the community that a security release is coming, to give users time to prepare their systems for the update. Notifications can include Discord messages, tweets, and emails to partners and validators.
- 24 hours following this notification, the fixes are applied publicly and new releases are issued.
- Once releases are available for Terra, we notify the community, again, through the same channels as above. We also publish a Security Advisory on Github and publish the CVE, as long as neither the Security Advisory nor the CVE include any information on how to exploit these vulnerabilities beyond what information is already available in the patch itself.
- Once the community is notified, we will pay out any relevant bug bounties to submitters.
- One week after the releases go out, we will publish a post with further details on the vulnerability as well as our response to it.

This process can take some time. Every effort will be made to handle the bug in as timely a manner as possible. However, it's important that we follow the process described above to ensure that disclosures are handled consistently and to keep Terra and the projects running on it secure.
