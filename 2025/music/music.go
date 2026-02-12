package music

var U2 = []album{
   {
      name: "war",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_nuez2OLRZURWEwEonP_eKZIjKKP_o93e0"},
      },
      track: []track{
         {
            number: 1,
            name: "sunday bloody sunday",
            link: []link{
               {url: "youtube.com/watch?v=d8FYJ9qvegI"},
            },
         },
      },
   },
}

var EsbenAndTheWitch = []album{
   {
      name: "Wash the Sins Not Only the Face",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_mWTxrHEOalIe-BLQrIs_kx-plYsyXysdM"},
      },
      track: []track{
         {
            number: 2,
            name: "slow wave",
            link: []link{
               {url: "youtube.com/watch?v=2LYLeVar5OA"},
            },
         },
         {
            number: 6,
            name: "yellow wood",
            link: []link{
               {url: "youtube.com/watch?v=-VHsbhM6o9A"},
            },
         },
         {
            number: 9,
            name: "the fall of glorieta mountain",
            link: []link{
               {url: "youtube.com/watch?v=qqXsyk5ZmCo"},
            },
         },
      },
   },
}

var KellyMoran = []album{
   {
      name: "helix (edit)",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_kB2WFbIR3N8VhQWq0G94nUUoa1275EcMU"},
      },
      track: []track{
         {
            number: 1,
            name:   "helix (edit)",
            link: []link{
               {url: "youtube.com/watch?v=tdeZ4ecFaxE"},
            },
         },
      },
   },
   {
      name: "moves in the field",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_mI7sh2DxPKFiColG-vykv9e74yqdKQuMI"},
      },
      track: []track{
         {
            number: 4,
            name:   "Dancer Polynomials",
            link: []link{
               {url: "youtube.com/watch?v=E6LxCW7cl4E"},
            },
         },
         {
            number: 5,
            name:   "sodalis",
            link: []link{
               {url: "youtube.com/watch?v=VBlhTlz1LXA"},
            },
         },
         {
            number: 6,
            name:   "leitmotif",
            link: []link{
               {url: "youtube.com/watch?v=o63kYG-j6Lg"},
            },
         },
         {
            number: 7,
            name:   "it's okay to disappear",
            link: []link{
               {url: "youtube.com/watch?v=exedtrCI3fg"},
            },
         },
      },
   },
   {
      name: "ultraviolet",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_mQz8ftbRzC-chKW0YSJ6yyPmgdrHv12CI"},
      },
      track: []track{
         {
            number: 1,
            name:   "autowave",
            link: []link{
               {url: "youtube.com/watch?v=vD1m4flw_Ko"},
            },
         },
      },
   },
}

var Bibio = []album{
   {
      name: "Phantom Brickworks",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_n-QU7H56CP_Cxg7Qg-ZKsHfUiFkh518Bw"},
      },
      track: []track{
         {
            number: 5,
            name:   "CAPEL CELYN",
            link: []link{
               {url: "youtube.com/watch?v=CJsBL-6Ixb8"},
               {url: "youtube.com/watch?v=nCuoUr5c_rI"},
            },
         },
         {
            number: 7,
            name:   "Bibio Ivy Charcoal",
            link: []link{
               {url: "youtube.com/watch?v=w7y8-lXjGAw"},
            },
         },
      },
   },
   {
      name: "Phantom Brickworks (IV & V)",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_m9SC1PouugkPkMBuKHH1J2TWPb0Bn1Mno"},
      },
      track: []track{
         {
            number: 2,
            name:   "PHANTOM BRICKWORKS V",
            link: []link{
               {url: "youtube.com/watch?v=akFKIFhz5iU"},
            },
         },
      },
   },
}

var MkGee = []album{
   {
      name: "rockman",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_meKfaX4xfVCc1wUZMYBKx0bzlGPYJUeTg"},
      },
      track: []track{
         {
            number: 1,
            name: "rockman",
            link: []link{
               {url: "youtube.com/watch?v=pULkm-3b_1M"},
            },
         },
      },
   },
   {
      name: "Two Star & the Dream Police",
      link: []link{
         {url: "youtube.com/playlist?list=OLAK5uy_lXnYM2XCJUqRw7PgY36hjIVKEOCMkkr-I"},
      },
      track: []track{
         {
            number: 3,
            name: "are you looking up",
            link: []link{
               {url: "youtube.com/watch?v=z0pzzkp85-Q"},
            },
         },
         {
            number: 11,
            name: "little bit more",
            link: []link{
               {url: "youtube.com/watch?v=NncelKQ6Hvw"},
            },
         },
         {
            number: 12,
            name: "dream police",
            link: []link{
               {url: "youtube.com/watch?v=3bLYswnD8yk"},
            },
         },
      },
   },
}

type album struct {
   name  string
   link  []link
   track []track
}

type track struct {
   number int
   name   string
   link   []link
}

type link struct {
   text string
   url  string
}
