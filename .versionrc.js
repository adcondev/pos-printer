module.exports = {
    // Solo mostrar lo importante en el changelog
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance"},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"},
        {type: "test", section: "✅ Tests"},
        {type: "ci", section: "🤖 Continuous Integration"},
        {type: "build", section: "🏗️ Build System"},
        {type: "style", section: "🎨 Styles"},
        {type: "refactor", section: "♻️ Code Refactoring"},
        {type: "chore", section: "🧹 Chores"},
        {type: "docs", section: "📝 Documentation"},
    ],

    // Configuración de GitHub
    commitUrlFormat: "https://github.com/adcondev/pos-printer/commit/{{hash}}",
    compareUrlFormat: "https://github.com/adcondev/pos-printer/compare/{{previousTag}}...{{currentTag}}",

    // Skip CI en commits de release
    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]"
};