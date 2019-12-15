import request from '@/utils/request'

export function addServer(data) {
  return request({
    url: '/redisServers',
    method: 'post',
    data
  })
}

export function updateServer(data) {
  return request({
    url: `/redisServers/${data.id}`,
    method: 'put',
    data
  })
}

export function listServer() {
  return request({
    url: `/redisServers`,
    method: 'get'
  })
}

export function removeServer(id) {
  return request({
    url: `/redisServers/${id}`,
    method: 'delete'
  })
}
