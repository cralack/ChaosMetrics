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

export const getHeroData = (name, loc, mode, version) => {
  return service({
    url: '/hero',
    method: 'get',
    params: {
      name,
      loc,
      mode,
      version,
    }
  })
}

export const getHeroDetail = (name, version, lang) => {
  return service({
    url: '/hero_detail',
    method: 'get',
    params: {
      name,
      version,
      lang
    }
  })
}

export const getPerks = (version, lang) => {
  return service({
    url: 'perks',
    method: 'get',
    params: {
      version,
      lang
    }
  })
}
