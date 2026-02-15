/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                // Semantic Palette
                bg: {
                    primary: 'var(--bg-primary)',
                    secondary: 'var(--bg-secondary)',
                    tertiary: 'var(--bg-tertiary)',
                },
                text: {
                    primary: 'var(--text-primary)',
                    secondary: 'var(--text-secondary)',
                    tertiary: 'var(--text-tertiary)',
                },
                border: {
                    primary: 'var(--border-primary)',
                    secondary: 'var(--border-secondary)',
                },
                accent: {
                    primary: 'var(--accent-primary)',
                    secondary: 'var(--accent-secondary)',
                },
                // Legacy support / specific overrides if needed
                vscode: {
                    bg: 'var(--bg-primary)',
                    text: 'var(--text-primary)',
                    primary: 'var(--accent-primary)',
                    sidebar: 'var(--bg-secondary)',
                    border: 'var(--border-primary)',
                }
            },
            fontFamily: {
                sans: ['Inter', 'system-ui', 'sans-serif'],
                mono: ['Fira Code', 'monospace'],
            },
            boxShadow: {
                'glass': '0 8px 32px 0 rgba(31, 38, 135, 0.15)',
                'glow': '0 0 15px rgba(130, 170, 255, 0.3)',
            },
        },
    },
    plugins: [],
}
