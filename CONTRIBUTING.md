# Contributing

We are always happy about support for our project. If you dicide to help us spreading awareness for privacy, just start with the following steps: 
1) Fork the project
2) Create a new branch
3) Commit your changes
4) Open a pull request to `master` 
  
(More about this in ["Branching & Naming"](#branching--naming))

For your Git Commit Messages, please orientate on the guidance in the following [article](https://chris.beams.io/posts/git-commit/):

> - Limit the subject line to 50 characters
> - Capitalize the subject line
> - Use the imperative mood in the subject line
> - Use the body to explain what and why, less how

### Pull Requests

Every new feature must be developed on a feature branch and merged into master. **Please do not push directly to master!** We also provide a [Pull Request Template](https://github.com/codeuniversity/smag-mvp/blob/master/.github/pull_request_template.md) for additional guidance. In any case, the Pull Request has to be reviewed and approved by at least one other developer before merging. Please make sure to [reference the associated issue(s)](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords) in the pull request.

## Branching & Naming

Next to `feature/<feature-name>` and `fix/<bug-name>` branches, we also have a `master` and a `production` branch. `master` is our development branch were new code is merged into first - release versions for roll-out will be merged into `production`.

For the naming of software components, please orientate on the existing components and folders. Especially if you build something specific to one of the platforms (e.g `insta(gram)`, `twitter`, ...), please make sure to use the regarding prefix for folder-names. Else, please stick to formatting conventions for [Golang](https://golang.org/doc/effective_go.html), [Python](https://www.python.org/dev/peps/pep-0008/) and [React.js](https://hackernoon.com/structuring-projects-and-naming-components-in-react-1261b6e18d76).

## Task Management

For our task management, we are using the [ZenHub GitHub Extension](https://www.zenhub.com/extension) which integrates a project board into GitHub. After installing the extension and reloading your browser, you will be able to see an addional `ZenHub` Tab in our repo. In there, you can see our current tasks `"In Progress"` and upcomming tasks `"ToDo"` of the current release we are working on. All tasks are represented as GitHub issues as well, so you might want to [create an own GitHub issue](https://github.com/codeuniversity/smag-mvp/issues/new/choose) for the beginning.

If you have any questions or want to get more involved in the project, feel free to approach the team via: [socialrecord-project[at]code.berlin](socialrecord-project@code.berlin).
