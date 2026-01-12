/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'vscode-bg': '#1e1e1e',
                'vscode-text': '#d4d4d4',
                'vscode-primary': '#007acc',
                'vscode-sidebar': '#252526',
                'vscode-border': '#3e3e42',
            },
        },
    },
    plugins: [],
}
