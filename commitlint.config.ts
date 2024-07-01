import type { UserConfig } from "@commitlint/types";
import { RuleConfigSeverity } from "@commitlint/types";

const Configuration: UserConfig = {
  extends: ["@commitlint/config-conventional"],
  formatter: "@commitlint/format",
  rules: {
    "header-max-length": [
      RuleConfigSeverity.Error,
      "always",
      Infinity,
    ] as const,
    "body-max-length": [RuleConfigSeverity.Error, "always", Infinity] as const,
    "body-max-line-length": [
      RuleConfigSeverity.Error,
      "always",
      Infinity,
    ] as const,
    "type-enum": [
      RuleConfigSeverity.Error,
      "always",
      [
        "feat",
        "fix",
        "perf",
        "revert",
        "docs",
        "chore",
        "ref",
        "test",
        "ci",
        "deps",
        "build",
      ],
    ],
  },
  prompt: {
    questions: {
      type: {
        description: "Select the type of change that you're committing",
        enum: {
          feat: {
            description: "A new feature",
            title: "Features",
            emoji: "🆕",
          },
          fix: {
            description: "A bug fix",
            title: "Bug Fixes",
            emoji: "🪲",
          },
          perf: {
            description: "A code change that improves performance",
            title: "Performance Improvements",
            emoji: "🚀",
          },
          revert: {
            description: "Reverts a previous commit",
            title: "Reverts",
            emoji: "⏪",
          },
          docs: {
            description: "Documentation only changes",
            title: "Documentation",
            emoji: "📑",
          },
          chore: {
            description: "Other changes that don't modify src or test files",
            title: "Chores",
            emoji: "⚓",
          },
          ref: {
            description:
              "A code change that neither fixes a bug nor adds a feature",
            title: "Code Refactoring",
            emoji: "🧹",
          },
          test: {
            description: "Adding missing tests or correcting existing tests",
            title: "Tests",
            emoji: "⌛",
          },
          ci: {
            description:
              "Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)",
            title: "Continuous Integrations",
            emoji: "♾️",
          },
          deps: {
            description: "Dependencies",
            title: "Dependencies",
            emoji: "🛠️",
          },
          build: {
            description: "Build",
            title: "Build",
            emoji: "🏗️",
          },
        },
      },
    },
  },
};

export default Configuration;
