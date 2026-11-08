# CLAUDE.md

## Pull requests come with a linked issue

Every pull request you open in this repo (`CerealKiller97/linkwarden-cli`) must close a GitHub issue.

1. **Find or create the issue before you open the PR.**
   - If an open issue already describes the change, reuse it. Look with `gh issue list --search "<keywords>"`.
   - Otherwise create one with `gh issue create`. Describe the problem or the need, not the fix:
     - **Title:** a plain summary of the problem, with no commit-type prefix. For example, "Spinner label disappears before it can be read".
     - **Body:** what is wrong or missing, how to see it, and why it matters.
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
