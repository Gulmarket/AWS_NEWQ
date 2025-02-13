const { join } = require('node:path')
module.exports = {
  root: true,
  plugins: ['@typescript-eslint'],
  parser: '@typescript-eslint/parser',
  extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended'],
  parserOptions: { ecmaVersion: 2018, sourceType: 'module' },
  rules: {
    'no-unused-vars': 'off'
  }
}
