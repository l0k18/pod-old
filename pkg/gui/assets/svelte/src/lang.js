// defining a locale dictionary
dictionary.set({
  
})

locale.set(
  getClientLocale({
    navigator: true,
    hash: 'lang',
    fallback: 'en',
  }),
)

locale.subscribe(l => {
  console.log('locale change', l)
})