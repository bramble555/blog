import axios from './axios'

export function postArticleDigg(sn) {
  return axios.post('/articles/digg', { sn })
}

export function postCommentDigg(sn) {
  return axios.post('/comments/digg', { sn })
}
