import request from '@/utils/request'

export function getDbCount(params) {
  return request({
    url: '/redis/dbCount',
    method: 'get',
    params
  })
}

export function keyTreeNodes(params) {
  return request({
    url: '/redis/keys/treeNodes',
    method: 'get',
    params
  })
}

export function keySummary(params) {
  return request({
    url: '/redis/keys/summary',
    method: 'get',
    params
  })
}

export function getValue(params) {
  return request({
    url: 'redis/keys/value',
    method: 'get',
    params
  })
}
