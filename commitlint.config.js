module.exports = {
    extends: ['@commitlint/config-conventional'],
    ignores: [(message) => /^bumps \[.+]\(.+\) from .+ to .+\.$/m.test(message)],
};
