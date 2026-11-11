# CLAUDE.md

## Pull requests come with a linked issue

Every pull request you open in this repo (`CerealKiller97/linkwarden-cli`) must close a GitHub issue.

1. **Find or create the issue before you open the PR.**
   - If an open issue already describes the change, reuse it. Look with `gh issue list --search "<keywords>"`.
   - Otherwise create one with `gh issue create`. Describe the problem or the need, not the fix:
     - **Title:** a plain summary of the problem, with no commit-type prefix. For example, "Spinner label disappears before it can be read".
     - **Body:** what is wrong or missing, how to see it, and why it matters. Follow the matching issue form (see [Issue and PR templates](#issue-and-pr-templates)).
     - **Label:** match the change type:

       | Commit type | Label |
       |---|---|
       | `fix` | `bug` |
       | `feat` | `enhancement` |
       | `docs` | `documentation` |
       | anything else (`refactor`, `test`, `chore`, `style`, …) | none |

2. **Link the issue in the PR body.** Put `Closes #<number>` on its own line near the top. GitHub then links the two and closes the issue when the PR merges. Use one closing line per issue.

3. **When a PR covers more than one change** (for example, several commits of different types), open one issue per change and add a `Closes #<n>` line for each.

4. **For a PR that already exists without an issue**, create the issue, then add the `Closes #<n>` line with `gh pr edit <pr> --body-file -`. Keep the rest of the body as it was.

5. **Report both links** to the user, the issue's and the PR's.

Skip the issue only when the user says to for that PR.

## Issue and PR templates

The repo has issue forms and a PR template in `.github/`. Blank issues are turned off, so people on GitHub always pick a form. `gh issue create` can't fill an issue form, so when you create an issue or PR from the CLI, write the body yourself in the same shape.

| File | Use for | Label |
|---|---|---|
| `.github/ISSUE_TEMPLATE/bug_report.yml` | `fix` changes | `bug` |
| `.github/ISSUE_TEMPLATE/feature_request.yml` | `feat` changes | `enhancement` |
| `.github/ISSUE_TEMPLATE/config.yml` | contact links (security, upstream Linkwarden) | n/a |
| `.github/pull_request_template.md` | every PR | n/a |

- **Issue body:** read the matching form and use each field's `label` as a `### ` heading, in the form's order, the way GitHub renders a submitted form. Fill the required fields. For an optional field you have nothing for, write `_No response_`. Issues of other types (`docs`, `refactor`, …) have no form: write a short problem description.
- **PR body:** read `.github/pull_request_template.md` and fill it in. Keep its sections and order:
  - Replace the top `Closes #` line with the real `Closes #<n>` line(s) and drop the HTML comment under it.
  - Tick every PR type that applies. List changes under the right Added / Changed / Removed / Fixed / Security / Deprecated heading, and leave `- None` under the rest.
  - Tick one tests box. If no tests changed, say why.
  - End the body with the attribution line, after "Further comments".
- **Security problems** never go in a public issue. Point the user to the private advisory link in `config.yml`.
