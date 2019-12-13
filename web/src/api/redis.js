import request from '@/utils/request'

export function getKeys(params) {
  return request({
    url: '/table/list',
    method: 'get',
    params
  })
}
