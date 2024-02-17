module.exports = {
  // ...
  extends: [
    // ...
    'plugin:astro/recommended',
    'plugin:@typescript-eslint/recommended', // Add TypeScript recommendations
    'plugin:react/recommended', // Add React recommendations (if you're using React)
    'plugin:react-hooks/recommended', // Add React hooks recommendations (if you're using React hooks)
    'prettier', // Add Prettier plugin
  ],
  settings: {
    react: {
      version: 'detect', // Automatically detect the React version
    },
  },
  overrides: [
    {
      // Define the configuration for `.astro` file.
      files: ['*.astro'],
      // Allows Astro components to be parsed.
      parser: 'astro-eslint-parser',
      // Parse the script in `.astro` as TypeScript by adding the following configuration.
      // It's the setting you need when using TypeScript.
      parserOptions: {
        parser: '@typescript-eslint/parser',
        extraFileExtensions: ['.astro'],
      },
      rules: {
        'react/react-in-jsx-scope': 'off', // Turn off the requirement for React in scope
        'react/no-unknown-property': ['error', { ignore: ['class'] }], // Allow 'class' property
        'react/jsx-key': 'off', // Disable the jsx-key rule for Astro files
        'react/no-unknown-property': 'off', // Disable the rule for unknown React properties in Astro files
      },
    },
    {
      files: ['*.ts', '*.tsx'], // Target TypeScript and TSX files
      parser: '@typescript-eslint/parser', // Use TypeScript parser
      parserOptions: {
        project: './tsconfig.json', // Path to your tsconfig.json
      },
      rules: {
        // TypeScript and React specific rules
        '@typescript-eslint/explicit-function-return-type': 'off',
        'react/react-in-jsx-scope': 'off',
        // Other TypeScript/React specific rules...
      },
    },
  ],
};
