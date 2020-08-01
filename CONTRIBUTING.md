# Contributing to the OneFuse Terraform Provider

Thank you for checking out the contributing docs!
The Terraform Provider for OneFuse is maintained by CloudBolt Software, Inc, but we would love to have you help out by providing your feedback and improvements.
This document outlines guidelines for contributing to the Terraform Provider for OneFuse.

## Making Issues

If you would like to share a bug, or a potential improvement, please make an Issue on Github.

Please include relevant information like:

* Your use-case for this feature
* Your system information including which versions of OneFuse, the provider, and Terraform you are using.
* Anything else that will help us scope out the work

## Making Pull Requests

If you have a bugfix or feature you would like to contribute, please make a Pull Request on Github.

1. Fork the repository to your own GitHub organization.
2. Add your work to that fork.
3. Add automated tests, where possible.
4. Ensure that automated tests continue to pass.
5. Run `make fmt` on your code and commit any changes it makes.
6. Make a Pull Request to the main repository.

For large feature work, you should make an Issue before starting work on the feature, as a place to discuss implementation details.

For obvious bugs, an issue is not required, but would be helpful to have.

Take pride in your work!
Include the following in you Pull Request description:

* What the changes have been made.
* Questions or concerns (if any) you have about the work.
* How to test the changes (manually or automated).
* Future work that could be done.
* References to related Issues/Pull Requests.

## What to expect

Somebody from CloudBolt will try to review your Issue/Pull Request within a few days of posting.

For Issues, that contributor will put the issue in the Development Backlog and try to provide an estimate about when the work will be done.
If the Issue warrants a discussion, they will respond and coordinate soliciting other feedback from within CloudBolt.

For Pull Requests, that contributor will provide feedback such as:

* Changes to the implementation.
* Changes to the testing.
* Changes to the documentation.

If necessary (and if you are OK with it) a CloudBolt developer may also add changes to the Pull Request themselves.

If you do not get a response within a week of posting, feel free to ask in the *#terraform-provider* channel of the [CloudBolt Users Slack](cloudbolt-users.slack.com).

## Building Releases

To cut a release of the provider, follow these steps:

1. Update the string in `VERSION`.
    * This should follow [semver](https://semver.org/) and follow the form `X.Y.Z` where
        * `X` increases for breaking changes
        * `Y` increases for new features
        * `Z` increases for bugfixes
    * The updated `VERSION` file should be committed.
1. Update the CHANGELOG to capture all changes since the last release.
    * Ideally the CHANGELOG is updated as changes are merged to the repo and this step is double-checking.
1. Create a tag for this version.
1. Push both the tag and the version to Github.
1. Run `make release` to generate the release files.
1. Draft a release in Github.
    * Include the binaries and `.sha256` files generated in the previous step.
    * For the body of the release, incude the CHANGELOG information for this release.
    * Have fun with it. Include a some relevant emoji.

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic address, without explicit permission
* Other conduct which could reasonably be considered inappropriate in a professional setting

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable behavior and are expected to take appropriate and fair corrective action in response to any instances of unacceptable behavior.

Project maintainers have the right and responsibility to remove, edit, or reject comments, commits, code, wiki edits, issues, and other contributions that are not aligned to this Code of Conduct, or to ban temporarily or permanently any contributor for other behaviors that they deem inappropriate, threatening, offensive, or harmful.

### Scope

This Code of Conduct applies both within project spaces and in public spaces when an individual is representing the project or its community. Examples of representing a project or community include using an official project e-mail address, posting via an official social media account, or acting as an appointed representative at an online or offline event. Representation of a project may be further defined and clarified by project maintainers.

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team at the `terraform-dev[at]cloudbolt.io`.
All complaints will be reviewed and investigated and will result in a response that is deemed necessary and appropriate to the circumstances.
The project team is obligated to maintain confidentiality with regard to the reporter of an incident.
Further details of specific enforcement policies may be posted separately.

Project maintainers who do not follow or enforce the Code of Conduct in good faith may face temporary or permanent repercussions as determined by other members of the project's leadership.

### Attribution

This Code of Conduct is adapted from the Contributor Covenant, version 1.4, available at http://contributor-covenant.org/version/1/4
