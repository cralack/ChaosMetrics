import service from '@/api/axios'

export const getGameVersion = () => {
  return service({
    url: '/gameversion',
    method: 'get'
  })
}

export const getClassicChampionRankBrief = (loc, version) => {
  return service({
    url: '/CLASSIC',
    method: 'get',
    params: {
      loc,
      version,
    }
  })
}
export const getARAMChampionRankBrief = (loc, version) => {
  return service({
    url: '/ARAM',
    method: 'get',
    params: {
      loc,
      version,
    }
  })
}
