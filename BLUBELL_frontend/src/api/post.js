import request from './request'

// 获取帖子列表 (传统分页)
export function getPostList(params) {
  return request({
    url: '/posts',
    method: 'get',
    params
  })
}

// 根据时间或分数获取帖子列表 (高级分页)
export function getPostList2(params) {
  return request({
    url: '/posts2',
    method: 'get',
    params
  })
}

// 获取帖子详情
export function getPostDetail(id) {
  return request({
    url: `/post/${id}`,
    method: 'get'
  })
}

// 创建帖子
export function createPost(data) {
  return request({
    url: '/post',
    method: 'post',
    data
  })
}

// 投票
export function votePost(data) {
  return request({
    url: '/vote',
    method: 'post',
    data
  })
}

// 获取社区列表
export function getCommunityList() {
  return request({
    url: '/community',
    method: 'get'
  })
}

// 获取社区详情
export function getCommunityDetail(id) {
  return request({
    url: `/community/${id}`,
    method: 'get'
  })
}
