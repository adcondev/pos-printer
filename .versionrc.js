module.exports = {
    // Updated types with library-specific sections
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance"},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"},
        {type: "test", section: "✅ Tests"},
        {type: "ci", section: "🤖 Continuous Integration"},
        {type: "build", section: "🏗️ Build System"},
        {type: "refactor", section: "♻️ Code Refactoring"},
        {type: "docs", section: "📝 Documentation"},
        {type: "style", section: "🎨 Code Style"},
        {type: "chore", hidden: true}
    ],

    // GitHub configuration
    commitUrlFormat: "https://github.com/adcondev/pos-printer/commit/{{hash}}",
    compareUrlFormat: "https://github.com/adcondev/pos-printer/compare/{{previousTag}}...{{currentTag}}",
    userUrlFormat: "https://github.com/{{user}}",

    // Skip CI on release commits
    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]",

    // Custom header for CHANGELOG
    header: "# Changelog\n\nAll notable changes to the POS Printer library will be documented in this file.\n"
};