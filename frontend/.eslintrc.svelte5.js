/**
 * ESLint Configuration for Svelte 5 + Vite Projects
 *
 * This config adds rules to catch common issues identified in Epic 3 retrospective.
 *
 * Usage: Add `"extends": ["./.eslintrc.svelte5.js"]` to your main .eslintrc.js
 * Or run: eslint --config .eslintrc.svelte5.js src/
 */

module.exports = {
    env: {
        browser: true,
        es2022: true,
        node: true
    },
    parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module'
    },
    plugins: ['svelte'],

    // Global rules
    rules: {
        // Consistent spacing
        'array-bracket-spacing': ['error', 'never'],
        'object-curly-spacing': ['error', 'always']
    },

    overrides: [
        {
            files: ['*.svelte'],
            parser: 'svelte-eslint-parser',
            rules: {
                // Prevent direct mutation patterns
                'no-param-reassign': ['error', {
                    props: true,
                    ignorePropertyModificationsFor: ['state', 'store', 'obj']
                }],

                // Prefer template literals
                'prefer-template': 'error',

                // Require accessibility attributes (Svelte-specific)
                'svelte/aria-attributes': 'error',
                'svelte/role-has-required-aria-attributes': 'error',
                'svelte/no-noninteractive-tabindex': 'warn',

                // No console.log in production code
                'no-console': ['warn', { allow: ['warn', 'error'] }]
            }
        },
        {
            files: ['*.js'],
            rules: {
                // Prevent direct mutation
                'no-param-reassign': ['error', {
                    props: true,
                    ignorePropertyModificationsFor: ['state', 'store']
                }]
            }
        }
    ]
};

/*
 * Rules Description:
 *
 * 1. no-param-reassign
 *    - Catches direct mutation like `items.push(x)` where items is a param
 *    - Forces creation of new references
 *
 * 2. prefer-template
 *    - Encourages template literals over string concatenation
 *
 * 3. svelte/aria-attributes
 *    - Requires aria-* attributes on interactive elements (Svelte-specific)
 *
 * 4. no-console
 *    - Warns about console.log (allow warn/error for logging)
 *
 * Note: Hardcoded URL detection is handled via CHECKLIST.md and code review,
 * not ESLint (no-restricted-globals doesn't work for string patterns).
 *
 * Installation:
 *   npm install -D eslint eslint-plugin-svelte svelte-eslint-parser @typescript-eslint/parser
 *
 * Add to .eslintrc.js:
 *   module.exports = {
 *       extends: ['./.eslintrc.svelte5.js'],
 *       // your other config
 *   };
 */
