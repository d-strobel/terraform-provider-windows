To import a local group member, you need both the SID (Security Identifier) of the group and the SID of the group member.
These two SIDs should be concatenated with `/member/` in between.

Format:
`<Group SID>/member/<Member SID>`
