const domains = ["tempmail.io", "disposable.com", "throwaway.net", "quickmail.org", "tempinbox.me"]

const prefixes = ["temp", "disposable", "throwaway", "quick", "random", "secure", "private", "anon"]

export function generateRandomEmail(): string {
  const randomPrefix = prefixes[Math.floor(Math.random() * prefixes.length)]
  const randomDomain = domains[Math.floor(Math.random() * domains.length)]
  const randomNumber = Math.floor(Math.random() * 10000)

  return `${randomPrefix}${randomNumber}@${randomDomain}`
}

