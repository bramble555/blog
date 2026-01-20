export const formatUrl = (url) => {
  if (!url) return ''
  const u = String(url).trim()
  if (!u) return ''
  if (u.startsWith('http')) return u
  if (u.startsWith('/')) return 'http://localhost:8080' + u
  return 'http://localhost:8080/' + u
}
