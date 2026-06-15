# Core Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 實作 YaoNomad 三硬幣卜卦核心：共用資料（hexagrams.json）、Go TUI 核心、Web 純 JS 核心，均含完整測試。

**Architecture:** core-spec/hexagrams.json 為唯一資料正本，兩端各持副本。Go 核心位於 tui-go/internal/model（types/casting/hexagrams），Web 核心位於 web/core（純 JS ES Modules）。查卦使用線性掃描，隨機源注入以利測試。

**Tech Stack:** Go 1.21+（標準庫 math/rand、encoding/json、testing），Node.js 18+（node:test、node:assert、ES Modules），無外部依賴。

**JS 弱型別注意事項：** 全面使用 `===`，禁用 `==`。所有函式以 JSDoc 標註型別。布林判斷用 `l.yin === true` 或直接 `l.yin`，不依賴 truthy/falsy 轉換。

---

## 檔案結構

| 動作 | 路徑 |
|------|------|
| 建立 | `core-spec/README.md` |
| 建立 | `core-spec/hexagrams.json` ← **正本** |
| 建立 | `tui-go/go.mod` |
| 建立 | `tui-go/internal/model/types.go` |
| 建立 | `tui-go/internal/model/hexagrams.go` |
| 建立 | `tui-go/internal/model/hexagrams_test.go` |
| 建立 | `tui-go/internal/model/casting.go` |
| 建立 | `tui-go/internal/model/casting_test.go` |
| 複製 | `tui-go/data/hexagrams.json` ← 從正本複製 |
| 建立 | `web/core/models.js` |
| 建立 | `web/core/hexagrams.js` |
| 建立 | `web/core/hexagrams.test.js` |
| 建立 | `web/core/casting.js` |
| 建立 | `web/core/casting.test.js` |
| 複製 | `web/data/hexagrams.json` ← 從正本複製 |

---

## Task 1：建立 core-spec 目錄與 hexagrams.json

**Files:**
- Create: `core-spec/README.md`
- Create: `core-spec/hexagrams.json`

- [ ] **Step 1：建立目錄**

```bash
mkdir -p core-spec
```

- [ ] **Step 2：建立 core-spec/README.md**

```bash
cat > core-spec/README.md << 'EOF'
# YaoNomad Core Spec

語言無關的卜卦演算法規格。詳見 docs/superpowers/specs/2026-06-15-core-design.md。

## 八卦三位碼（bit0=底爻，yang=1，yin=0）

| 卦 | 符號 | 碼 |
|----|------|----|
| 乾 | ☰ | 7 |
| 兌 | ☱ | 3 |
| 離 | ☲ | 5 |
| 震 | ☳ | 1 |
| 巽 | ☴ | 6 |
| 坎 | ☵ | 2 |
| 艮 | ☶ | 4 |
| 坤 | ☷ | 0 |
EOF
```

- [ ] **Step 3：建立 core-spec/hexagrams.json（完整 64 卦）**

寫入以下內容至 `core-spec/hexagrams.json`：

```json
[
  {
    "number": 1, "name": "乾", "upper": 7, "lower": 7,
    "judgment": "元亨利貞。",
    "lines": [
      "初九：潛龍，勿用。",
      "九二：見龍在田，利見大人。",
      "九三：君子終日乾乾，夕惕若，厲，無咎。",
      "九四：或躍在淵，無咎。",
      "九五：飛龍在天，利見大人。",
      "上九：亢龍有悔。"
    ]
  },
  {
    "number": 2, "name": "坤", "upper": 0, "lower": 0,
    "judgment": "元亨，利牝馬之貞。君子有攸往，先迷後得主。利西南得朋，東北喪朋。安貞吉。",
    "lines": [
      "初六：履霜，堅冰至。",
      "六二：直方大，不習無不利。",
      "六三：含章可貞，或從王事，無成有終。",
      "六四：括囊，無咎無譽。",
      "六五：黃裳，元吉。",
      "上六：龍戰于野，其血玄黃。"
    ]
  },
  {
    "number": 3, "name": "屯", "upper": 2, "lower": 1,
    "judgment": "元亨利貞。勿用有攸往。利建侯。",
    "lines": [
      "初九：磐桓，利居貞，利建侯。",
      "六二：屯如邅如，乘馬班如。匪寇婚媾，女子貞不字，十年乃字。",
      "六三：即鹿無虞，惟入于林中，君子幾不如舍，往吝。",
      "六四：乘馬班如，求婚媾，往吉，無不利。",
      "九五：屯其膏，小貞吉，大貞凶。",
      "上六：乘馬班如，泣血漣如。"
    ]
  },
  {
    "number": 4, "name": "蒙", "upper": 4, "lower": 2,
    "judgment": "亨。匪我求童蒙，童蒙求我。初筮告，再三瀆，瀆則不告。利貞。",
    "lines": [
      "初六：發蒙，利用刑人，用說桎梏，以往吝。",
      "九二：包蒙，吉。納婦吉，子克家。",
      "六三：勿用取女，見金夫，不有躬，無攸利。",
      "六四：困蒙，吝。",
      "六五：童蒙，吉。",
      "上九：擊蒙，不利為寇，利禦寇。"
    ]
  },
  {
    "number": 5, "name": "需", "upper": 2, "lower": 7,
    "judgment": "有孚，光亨，貞吉。利涉大川。",
    "lines": [
      "初九：需于郊，利用恆，無咎。",
      "九二：需于沙，小有言，終吉。",
      "九三：需于泥，致寇至。",
      "六四：需于血，出自穴。",
      "九五：需于酒食，貞吉。",
      "上六：入于穴，有不速之客三人來，敬之終吉。"
    ]
  },
  {
    "number": 6, "name": "訟", "upper": 7, "lower": 2,
    "judgment": "有孚，窒惕，中吉，終凶。利見大人，不利涉大川。",
    "lines": [
      "初六：不永所事，小有言，終吉。",
      "九二：不克訟，歸而逋，其邑人三百戶，無眚。",
      "六三：食舊德，貞厲，終吉。或從王事，無成。",
      "九四：不克訟，復即命渝，安貞吉。",
      "九五：訟，元吉。",
      "上九：或錫之鞶帶，終朝三褫之。"
    ]
  },
  {
    "number": 7, "name": "師", "upper": 0, "lower": 2,
    "judgment": "貞，丈人吉，無咎。",
    "lines": [
      "初六：師出以律，否臧凶。",
      "九二：在師中，吉無咎，王三錫命。",
      "六三：師或輿尸，凶。",
      "六四：師左次，無咎。",
      "六五：田有禽，利執言，無咎。長子帥師，弟子輿尸，貞凶。",
      "上六：大君有命，開國承家，小人勿用。"
    ]
  },
  {
    "number": 8, "name": "比", "upper": 2, "lower": 0,
    "judgment": "吉。原筮，元永貞，無咎。不寧方來，後夫凶。",
    "lines": [
      "初六：有孚比之，無咎。有孚盈缶，終來有它吉。",
      "六二：比之自內，貞吉。",
      "六三：比之匪人。",
      "六四：外比之，貞吉。",
      "九五：顯比，王用三驅，失前禽，邑人不誡，吉。",
      "上六：比之無首，凶。"
    ]
  },
  {
    "number": 9, "name": "小畜", "upper": 6, "lower": 7,
    "judgment": "亨。密雲不雨，自我西郊。",
    "lines": [
      "初九：復自道，何其咎，吉。",
      "九二：牽復，吉。",
      "九三：輿說輻，夫妻反目。",
      "六四：有孚，血去惕出，無咎。",
      "九五：有孚攣如，富以其鄰。",
      "上九：既雨既處，尚德載，婦貞厲。月幾望，君子征凶。"
    ]
  },
  {
    "number": 10, "name": "履", "upper": 7, "lower": 3,
    "judgment": "履虎尾，不咥人，亨。",
    "lines": [
      "初九：素履往，無咎。",
      "九二：履道坦坦，幽人貞吉。",
      "六三：眇能視，跛能履，履虎尾，咥人，凶。武人為于大君。",
      "九四：履虎尾，愬愬終吉。",
      "九五：夬履，貞厲。",
      "上九：視履考祥，其旋元吉。"
    ]
  },
  {
    "number": 11, "name": "泰", "upper": 0, "lower": 7,
    "judgment": "小往大來，吉亨。",
    "lines": [
      "初九：拔茅茹，以其彙，征吉。",
      "九二：包荒，用馮河，不遐遺，朋亡，得尚于中行。",
      "九三：無平不陂，無往不復，艱貞無咎。勿恤其孚，于食有福。",
      "六四：翩翩不富，以其鄰，不戒以孚。",
      "六五：帝乙歸妹，以祉元吉。",
      "上六：城復于隍，勿用師，自邑告命，貞吝。"
    ]
  },
  {
    "number": 12, "name": "否", "upper": 7, "lower": 0,
    "judgment": "否之匪人，不利君子貞，大往小來。",
    "lines": [
      "初六：拔茅茹，以其彙，貞吉亨。",
      "六二：包承，小人吉，大人否亨。",
      "六三：包羞。",
      "九四：有命無咎，疇離祉。",
      "九五：休否，大人吉。其亡其亡，繫于苞桑。",
      "上九：傾否，先否後喜。"
    ]
  },
  {
    "number": 13, "name": "同人", "upper": 7, "lower": 5,
    "judgment": "同人于野，亨。利涉大川，利君子貞。",
    "lines": [
      "初九：同人于門，無咎。",
      "六二：同人于宗，吝。",
      "九三：伏戎于莽，升其高陵，三歲不興。",
      "九四：乘其墉，弗克攻，吉。",
      "九五：同人先號咷而後笑，大師克相遇。",
      "上九：同人于郊，無悔。"
    ]
  },
  {
    "number": 14, "name": "大有", "upper": 5, "lower": 7,
    "judgment": "元亨。",
    "lines": [
      "初九：無交害，匪咎，艱則無咎。",
      "九二：大車以載，有攸往，無咎。",
      "九三：公用亨于天子，小人弗克。",
      "九四：匪其彭，無咎。",
      "六五：厥孚交如，威如，吉。",
      "上九：自天祐之，吉無不利。"
    ]
  },
  {
    "number": 15, "name": "謙", "upper": 0, "lower": 4,
    "judgment": "亨，君子有終。",
    "lines": [
      "初六：謙謙君子，用涉大川，吉。",
      "六二：鳴謙，貞吉。",
      "九三：勞謙君子，有終吉。",
      "六四：無不利，撝謙。",
      "六五：不富，以其鄰，利用侵伐，無不利。",
      "上六：鳴謙，利用行師，征邑國。"
    ]
  },
  {
    "number": 16, "name": "豫", "upper": 1, "lower": 0,
    "judgment": "利建侯行師。",
    "lines": [
      "初六：鳴豫，凶。",
      "六二：介于石，不終日，貞吉。",
      "六三：盱豫悔，遲有悔。",
      "九四：由豫，大有得，勿疑，朋盍簪。",
      "六五：貞疾，恆不死。",
      "上六：冥豫，成有渝，無咎。"
    ]
  },
  {
    "number": 17, "name": "隨", "upper": 3, "lower": 1,
    "judgment": "元亨利貞，無咎。",
    "lines": [
      "初九：官有渝，貞吉。出門交有功。",
      "六二：係小子，失丈夫。",
      "六三：係丈夫，失小子。隨有求得，利居貞。",
      "九四：隨有獲，貞凶。有孚在道，以明，何咎。",
      "九五：孚于嘉，吉。",
      "上六：拘係之，乃從維之，王用亨于西山。"
    ]
  },
  {
    "number": 18, "name": "蠱", "upper": 4, "lower": 6,
    "judgment": "元亨，利涉大川。先甲三日，後甲三日。",
    "lines": [
      "初六：幹父之蠱，有子考，無咎，厲終吉。",
      "九二：幹母之蠱，不可貞。",
      "九三：幹父之蠱，小有悔，無大咎。",
      "六四：裕父之蠱，往見吝。",
      "六五：幹父之蠱，用譽。",
      "上九：不事王侯，高尚其事。"
    ]
  },
  {
    "number": 19, "name": "臨", "upper": 0, "lower": 3,
    "judgment": "元亨利貞。至于八月有凶。",
    "lines": [
      "初九：咸臨，貞吉。",
      "九二：咸臨，吉，無不利。",
      "六三：甘臨，無攸利，既憂之，無咎。",
      "六四：至臨，無咎。",
      "六五：知臨，大君之宜，吉。",
      "上六：敦臨，吉，無咎。"
    ]
  },
  {
    "number": 20, "name": "觀", "upper": 6, "lower": 0,
    "judgment": "盥而不薦，有孚顒若。",
    "lines": [
      "初六：童觀，小人無咎，君子吝。",
      "六二：闚觀，利女貞。",
      "六三：觀我生，進退。",
      "六四：觀國之光，利用賓于王。",
      "九五：觀我生，君子無咎。",
      "上九：觀其生，君子無咎。"
    ]
  },
  {
    "number": 21, "name": "噬嗑", "upper": 5, "lower": 1,
    "judgment": "亨，利用獄。",
    "lines": [
      "初九：屨校滅趾，無咎。",
      "六二：噬膚，滅鼻，無咎。",
      "六三：噬腊肉，遇毒，小吝，無咎。",
      "九四：噬乾胏，得金矢，利艱貞，吉。",
      "六五：噬乾肉，得黃金，貞厲，無咎。",
      "上九：何校滅耳，凶。"
    ]
  },
  {
    "number": 22, "name": "賁", "upper": 4, "lower": 5,
    "judgment": "亨，小利有攸往。",
    "lines": [
      "初九：賁其趾，舍車而徒。",
      "六二：賁其須。",
      "九三：賁如濡如，永貞吉。",
      "六四：賁如皤如，白馬翰如，匪寇婚媾。",
      "六五：賁于丘園，束帛戔戔，吝，終吉。",
      "上九：白賁，無咎。"
    ]
  },
  {
    "number": 23, "name": "剝", "upper": 4, "lower": 0,
    "judgment": "不利有攸往。",
    "lines": [
      "初六：剝床以足，蔑貞凶。",
      "六二：剝床以辨，蔑貞凶。",
      "六三：剝之，無咎。",
      "六四：剝床以膚，凶。",
      "六五：貫魚，以宮人寵，無不利。",
      "上九：碩果不食，君子得輿，小人剝廬。"
    ]
  },
  {
    "number": 24, "name": "復", "upper": 0, "lower": 1,
    "judgment": "亨。出入無疾，朋來無咎。反復其道，七日來復。利有攸往。",
    "lines": [
      "初九：不遠復，無祗悔，元吉。",
      "六二：休復，吉。",
      "六三：頻復，厲無咎。",
      "六四：中行獨復。",
      "六五：敦復，無悔。",
      "上六：迷復，凶，有災眚。用行師，終有大敗，以其國君凶，至于十年不克征。"
    ]
  },
  {
    "number": 25, "name": "無妄", "upper": 7, "lower": 1,
    "judgment": "元亨利貞。其匪正有眚，不利有攸往。",
    "lines": [
      "初九：無妄，往吉。",
      "六二：不耕獲，不菑畬，則利有攸往。",
      "六三：無妄之災，或繫之牛，行人之得，邑人之災。",
      "九四：可貞，無咎。",
      "九五：無妄之疾，勿藥有喜。",
      "上九：無妄行，有眚，無攸利。"
    ]
  },
  {
    "number": 26, "name": "大畜", "upper": 4, "lower": 7,
    "judgment": "利貞，不家食吉，利涉大川。",
    "lines": [
      "初九：有厲，利已。",
      "九二：輿說輹。",
      "九三：良馬逐，利艱貞。曰閑輿衛，利有攸往。",
      "六四：童牛之牿，元吉。",
      "六五：豶豕之牙，吉。",
      "上九：何天之衢，亨。"
    ]
  },
  {
    "number": 27, "name": "頤", "upper": 4, "lower": 1,
    "judgment": "貞吉，觀頤，自求口實。",
    "lines": [
      "初九：舍爾靈龜，觀我朵頤，凶。",
      "六二：顛頤，拂經，于丘頤，征凶。",
      "六三：拂頤，貞凶，十年勿用，無攸利。",
      "六四：顛頤，吉。虎視眈眈，其欲逐逐，無咎。",
      "六五：拂經，居貞吉，不可涉大川。",
      "上九：由頤，厲吉，利涉大川。"
    ]
  },
  {
    "number": 28, "name": "大過", "upper": 3, "lower": 6,
    "judgment": "棟橈，利有攸往，亨。",
    "lines": [
      "初六：藉用白茅，無咎。",
      "九二：枯楊生稊，老夫得其女妻，無不利。",
      "九三：棟橈，凶。",
      "九四：棟隆，吉；有它，吝。",
      "九五：枯楊生華，老婦得其士夫，無咎無譽。",
      "上六：過涉滅頂，凶，無咎。"
    ]
  },
  {
    "number": 29, "name": "坎", "upper": 2, "lower": 2,
    "judgment": "習坎，有孚，維心亨，行有尚。",
    "lines": [
      "初六：習坎，入于坎窞，凶。",
      "九二：坎有險，求小得。",
      "六三：來之坎坎，險且枕，入于坎窞，勿用。",
      "六四：樽酒簋貳，用缶，納約自牖，終無咎。",
      "九五：坎不盈，祗既平，無咎。",
      "上六：係用徽纆，寘于叢棘，三歲不得，凶。"
    ]
  },
  {
    "number": 30, "name": "離", "upper": 5, "lower": 5,
    "judgment": "利貞，亨，畜牝牛吉。",
    "lines": [
      "初九：履錯然，敬之，無咎。",
      "六二：黃離，元吉。",
      "九三：日昃之離，不鼓缶而歌，則大耋之嗟，凶。",
      "九四：突如其來如，焚如，死如，棄如。",
      "六五：出涕沱若，戚嗟若，吉。",
      "上九：王用出征，有嘉折首，獲匪其醜，無咎。"
    ]
  },
  {
    "number": 31, "name": "咸", "upper": 3, "lower": 4,
    "judgment": "亨，利貞，取女吉。",
    "lines": [
      "初六：咸其拇。",
      "六二：咸其腓，凶，居吉。",
      "九三：咸其股，執其隨，往吝。",
      "九四：貞吉悔亡，憧憧往來，朋從爾思。",
      "九五：咸其脢，無悔。",
      "上六：咸其輔頰舌。"
    ]
  },
  {
    "number": 32, "name": "恆", "upper": 1, "lower": 6,
    "judgment": "亨，無咎，利貞，利有攸往。",
    "lines": [
      "初六：浚恆，貞凶，無攸利。",
      "九二：悔亡。",
      "九三：不恆其德，或承之羞，貞吝。",
      "九四：田無禽。",
      "六五：恆其德，貞，婦人吉，夫子凶。",
      "上六：振恆，凶。"
    ]
  },
  {
    "number": 33, "name": "遯", "upper": 7, "lower": 4,
    "judgment": "亨，小利貞。",
    "lines": [
      "初六：遯尾，厲，勿用有攸往。",
      "六二：執之用黃牛之革，莫之勝說。",
      "九三：係遯，有疾厲，畜臣妾吉。",
      "九四：好遯，君子吉，小人否。",
      "九五：嘉遯，貞吉。",
      "上九：肥遯，無不利。"
    ]
  },
  {
    "number": 34, "name": "大壯", "upper": 1, "lower": 7,
    "judgment": "利貞。",
    "lines": [
      "初九：壯于趾，征凶，有孚。",
      "九二：貞吉。",
      "九三：小人用壯，君子用罔，貞厲。羝羊觸藩，羸其角。",
      "九四：貞吉，悔亡，藩決不羸，壯于大輿之輹。",
      "六五：喪羊于易，無悔。",
      "上六：羝羊觸藩，不能退，不能遂，無攸利，艱則吉。"
    ]
  },
  {
    "number": 35, "name": "晉", "upper": 5, "lower": 0,
    "judgment": "康侯用錫馬蕃庶，晝日三接。",
    "lines": [
      "初六：晉如摧如，貞吉。罔孚，裕無咎。",
      "六二：晉如愁如，貞吉。受茲介福，于其王母。",
      "六三：眾允，悔亡。",
      "九四：晉如鼫鼠，貞厲。",
      "六五：悔亡，失得勿恤，往吉，無不利。",
      "上九：晉其角，維用伐邑，厲吉無咎，貞吝。"
    ]
  },
  {
    "number": 36, "name": "明夷", "upper": 0, "lower": 5,
    "judgment": "利艱貞。",
    "lines": [
      "初九：明夷于飛，垂其翼。君子于行，三日不食。有攸往，主人有言。",
      "六二：明夷，夷于左股，用拯馬壯，吉。",
      "九三：明夷于南狩，得其大首，不可疾貞。",
      "六四：入于左腹，獲明夷之心，于出門庭。",
      "六五：箕子之明夷，利貞。",
      "上六：不明晦，初登于天，後入于地。"
    ]
  },
  {
    "number": 37, "name": "家人", "upper": 6, "lower": 5,
    "judgment": "利女貞。",
    "lines": [
      "初九：閑有家，悔亡。",
      "六二：無攸遂，在中饋，貞吉。",
      "九三：家人嗃嗃，悔厲吉；婦子嘻嘻，終吝。",
      "六四：富家，大吉。",
      "九五：王假有家，勿恤，吉。",
      "上九：有孚威如，終吉。"
    ]
  },
  {
    "number": 38, "name": "睽", "upper": 5, "lower": 3,
    "judgment": "小事吉。",
    "lines": [
      "初九：悔亡，喪馬勿逐，自復。見惡人，無咎。",
      "九二：遇主于巷，無咎。",
      "六三：見輿曳，其牛掣，其人天且劓，無初有終。",
      "九四：睽孤，遇元夫，交孚，厲無咎。",
      "六五：悔亡，厥宗噬膚，往何咎。",
      "上九：睽孤，見豕負塗，載鬼一車，先張之弧，後說之弧，匪寇婚媾，往遇雨則吉。"
    ]
  },
  {
    "number": 39, "name": "蹇", "upper": 2, "lower": 4,
    "judgment": "利西南，不利東北。利見大人，貞吉。",
    "lines": [
      "初六：往蹇，來譽。",
      "六二：王臣蹇蹇，匪躬之故。",
      "九三：往蹇，來反。",
      "六四：往蹇，來連。",
      "九五：大蹇朋來。",
      "上六：往蹇，來碩，吉，利見大人。"
    ]
  },
  {
    "number": 40, "name": "解", "upper": 1, "lower": 2,
    "judgment": "利西南。無所往，其來復吉。有攸往，夙吉。",
    "lines": [
      "初六：無咎。",
      "九二：田獲三狐，得黃矢，貞吉。",
      "六三：負且乘，致寇至，貞吝。",
      "九四：解而拇，朋至斯孚。",
      "六五：君子維有解，吉，有孚于小人。",
      "上六：公用射隼，于高墉之上，獲之，無不利。"
    ]
  },
  {
    "number": 41, "name": "損", "upper": 4, "lower": 3,
    "judgment": "有孚，元吉，無咎，可貞，利有攸往。曷之用？二簋可用享。",
    "lines": [
      "初九：已事遄往，無咎，酌損之。",
      "九二：利貞，征凶，弗損益之。",
      "六三：三人行，則損一人；一人行，則得其友。",
      "六四：損其疾，使遄有喜，無咎。",
      "六五：或益之，十朋之龜弗克違，元吉。",
      "上九：弗損益之，無咎，貞吉，利有攸往，得臣無家。"
    ]
  },
  {
    "number": 42, "name": "益", "upper": 6, "lower": 1,
    "judgment": "利有攸往，利涉大川。",
    "lines": [
      "初九：利用為大作，元吉，無咎。",
      "六二：或益之，十朋之龜弗克違，永貞吉。王用享于帝，吉。",
      "六三：益之用凶事，無咎。有孚中行，告公用圭。",
      "六四：中行，告公從，利用為依遷國。",
      "九五：有孚惠心，勿問，元吉。有孚惠我德。",
      "上九：莫益之，或擊之，立心勿恆，凶。"
    ]
  },
  {
    "number": 43, "name": "夬", "upper": 3, "lower": 7,
    "judgment": "揚于王庭，孚號，有厲，告自邑，不利即戎，利有攸往。",
    "lines": [
      "初九：壯于前趾，往不勝，為咎。",
      "九二：惕號，莫夜有戎，勿恤。",
      "九三：壯于頄，有凶，君子夬夬，獨行遇雨，若濡有慍，無咎。",
      "九四：臀無膚，其行次且，牽羊悔亡，聞言不信。",
      "九五：莧陸夬夬，中行無咎。",
      "上六：無號，終有凶。"
    ]
  },
  {
    "number": 44, "name": "姤", "upper": 7, "lower": 6,
    "judgment": "女壯，勿用取女。",
    "lines": [
      "初六：繫于金柅，貞吉，有攸往，見凶，羸豕孚蹢躅。",
      "九二：包有魚，無咎，不利賓。",
      "九三：臀無膚，其行次且，厲，無大咎。",
      "九四：包無魚，起凶。",
      "九五：以杞包瓜，含章，有隕自天。",
      "上九：姤其角，吝，無咎。"
    ]
  },
  {
    "number": 45, "name": "萃", "upper": 3, "lower": 0,
    "judgment": "亨，王假有廟，利見大人，亨，利貞。用大牲吉，利有攸往。",
    "lines": [
      "初六：有孚不終，乃亂乃萃，若號，一握為笑，勿恤，往無咎。",
      "六二：引吉，無咎，孚乃利用禴。",
      "六三：萃如嗟如，無攸利，往無咎，小吝。",
      "九四：大吉，無咎。",
      "九五：萃有位，無咎，匪孚，元永貞，悔亡。",
      "上六：齎咨涕洟，無咎。"
    ]
  },
  {
    "number": 46, "name": "升", "upper": 0, "lower": 6,
    "judgment": "元亨，用見大人，勿恤，南征吉。",
    "lines": [
      "初六：允升，大吉。",
      "九二：孚乃利用禴，無咎。",
      "九三：升虛邑。",
      "六四：王用亨于岐山，吉，無咎。",
      "六五：貞吉，升階。",
      "上六：冥升，利于不息之貞。"
    ]
  },
  {
    "number": 47, "name": "困", "upper": 3, "lower": 2,
    "judgment": "亨，貞，大人吉，無咎，有言不信。",
    "lines": [
      "初六：臀困于株木，入于幽谷，三歲不覿。",
      "九二：困于酒食，朱紱方來，利用享祀，征凶，無咎。",
      "六三：困于石，據于蒺藜，入于其宮，不見其妻，凶。",
      "九四：來徐徐，困于金車，吝，有終。",
      "九五：劓刖，困于赤紱，乃徐有說，利用祭祀。",
      "上六：困于葛藟，于臲卼，曰動悔，有悔，征吉。"
    ]
  },
  {
    "number": 48, "name": "井", "upper": 2, "lower": 6,
    "judgment": "改邑不改井，無喪無得，往來井井。汔至亦未繘井，羸其瓶，凶。",
    "lines": [
      "初六：井泥不食，舊井無禽。",
      "九二：井谷射鮒，甕敝漏。",
      "九三：井渫不食，為我心惻，可用汲，王明並受其福。",
      "六四：井甃，無咎。",
      "九五：井冽，寒泉食。",
      "上六：井收勿幕，有孚元吉。"
    ]
  },
  {
    "number": 49, "name": "革", "upper": 3, "lower": 5,
    "judgment": "己日乃孚，元亨利貞，悔亡。",
    "lines": [
      "初九：鞏用黃牛之革。",
      "六二：己日乃革之，征吉，無咎。",
      "九三：征凶，貞厲，革言三就，有孚。",
      "九四：悔亡，有孚改命，吉。",
      "九五：大人虎變，未占有孚。",
      "上六：君子豹變，小人革面，征凶，居貞吉。"
    ]
  },
  {
    "number": 50, "name": "鼎", "upper": 5, "lower": 6,
    "judgment": "元吉，亨。",
    "lines": [
      "初六：鼎顛趾，利出否，得妾以其子，無咎。",
      "九二：鼎有實，我仇有疾，不我能即，吉。",
      "九三：鼎耳革，其行塞，雉膏不食，方雨虧悔，終吉。",
      "九四：鼎折足，覆公餗，其形渥，凶。",
      "六五：鼎黃耳金鉉，利貞。",
      "上九：鼎玉鉉，大吉，無不利。"
    ]
  },
  {
    "number": 51, "name": "震", "upper": 1, "lower": 1,
    "judgment": "亨。震來虩虩，笑言啞啞。震驚百里，不喪匕鬯。",
    "lines": [
      "初九：震來虩虩，後笑言啞啞，吉。",
      "六二：震來厲，億喪貝，躋于九陵，勿逐，七日得。",
      "六三：震蘇蘇，震行無眚。",
      "九四：震遂泥。",
      "六五：震往來厲，億無喪，有事。",
      "上六：震索索，視矍矍，征凶。震不于其躬，于其鄰，無咎，婚媾有言。"
    ]
  },
  {
    "number": 52, "name": "艮", "upper": 4, "lower": 4,
    "judgment": "艮其背，不獲其身，行其庭，不見其人，無咎。",
    "lines": [
      "初六：艮其趾，無咎，利永貞。",
      "六二：艮其腓，不拯其隨，其心不快。",
      "九三：艮其限，列其夤，厲薰心。",
      "六四：艮其身，無咎。",
      "六五：艮其輔，言有序，悔亡。",
      "上九：敦艮，吉。"
    ]
  },
  {
    "number": 53, "name": "漸", "upper": 6, "lower": 4,
    "judgment": "女歸吉，利貞。",
    "lines": [
      "初六：鴻漸于干，小子厲，有言，無咎。",
      "六二：鴻漸于磐，飲食衎衎，吉。",
      "九三：鴻漸于陸，夫征不復，婦孕不育，凶，利禦寇。",
      "六四：鴻漸于木，或得其桷，無咎。",
      "九五：鴻漸于陵，婦三歲不孕，終莫之勝，吉。",
      "上九：鴻漸于陸，其羽可用為儀，吉。"
    ]
  },
  {
    "number": 54, "name": "歸妹", "upper": 1, "lower": 3,
    "judgment": "征凶，無攸利。",
    "lines": [
      "初九：歸妹以娣，跛能履，征吉。",
      "九二：眇能視，利幽人之貞。",
      "六三：歸妹以須，反歸以娣。",
      "九四：歸妹愆期，遲歸有時。",
      "六五：帝乙歸妹，其君之袂，不如其娣之袂良，月幾望，吉。",
      "上六：女承筐無實，士刲羊無血，無攸利。"
    ]
  },
  {
    "number": 55, "name": "豐", "upper": 1, "lower": 5,
    "judgment": "亨，王假之，勿憂，宜日中。",
    "lines": [
      "初九：遇其配主，雖旬無咎，往有尚。",
      "六二：豐其蔀，日中見斗，往得疑疾，有孚發若，吉。",
      "九三：豐其沛，日中見沬，折其右肱，無咎。",
      "九四：豐其蔀，日中見斗，遇其夷主，吉。",
      "六五：來章，有慶譽，吉。",
      "上六：豐其屋，蔀其家，闚其戶，闃其無人，三歲不覿，凶。"
    ]
  },
  {
    "number": 56, "name": "旅", "upper": 5, "lower": 4,
    "judgment": "小亨，旅貞吉。",
    "lines": [
      "初六：旅瑣瑣，斯其所取災。",
      "六二：旅即次，懷其資，得童僕，貞。",
      "九三：旅焚其次，喪其童僕，貞厲。",
      "九四：旅于處，得其資斧，我心不快。",
      "六五：射雉，一矢亡，終以譽命。",
      "上九：鳥焚其巢，旅人先笑後號咷，喪牛于易，凶。"
    ]
  },
  {
    "number": 57, "name": "巽", "upper": 6, "lower": 6,
    "judgment": "小亨，利有攸往，利見大人。",
    "lines": [
      "初六：進退，利武人之貞。",
      "九二：巽在床下，用史巫紛若，吉，無咎。",
      "九三：頻巽，吝。",
      "六四：悔亡，田獲三品。",
      "九五：貞吉，悔亡，無不利。無初有終，先庚三日，後庚三日，吉。",
      "上九：巽在床下，喪其資斧，貞凶。"
    ]
  },
  {
    "number": 58, "name": "兌", "upper": 3, "lower": 3,
    "judgment": "亨，利貞。",
    "lines": [
      "初九：和兌，吉。",
      "九二：孚兌，吉，悔亡。",
      "六三：來兌，凶。",
      "九四：商兌，未寧，介疾有喜。",
      "九五：孚于剝，有厲。",
      "上六：引兌。"
    ]
  },
  {
    "number": 59, "name": "渙", "upper": 6, "lower": 2,
    "judgment": "亨，王假有廟，利涉大川，利貞。",
    "lines": [
      "初六：用拯馬壯，吉。",
      "九二：渙奔其機，悔亡。",
      "六三：渙其躬，無悔。",
      "六四：渙其群，元吉；渙有丘，匪夷所思。",
      "九五：渙汗其大號，渙王居，無咎。",
      "上九：渙其血，去逖出，無咎。"
    ]
  },
  {
    "number": 60, "name": "節", "upper": 2, "lower": 3,
    "judgment": "亨，苦節不可貞。",
    "lines": [
      "初九：不出戶庭，無咎。",
      "九二：不出門庭，凶。",
      "六三：不節若，則嗟若，無咎。",
      "六四：安節，亨。",
      "九五：甘節，吉，往有尚。",
      "上六：苦節，貞凶，悔亡。"
    ]
  },
  {
    "number": 61, "name": "中孚", "upper": 6, "lower": 3,
    "judgment": "豚魚吉，利涉大川，利貞。",
    "lines": [
      "初九：虞吉，有它不燕。",
      "九二：鳴鶴在陰，其子和之，我有好爵，吾與爾靡之。",
      "六三：得敵，或鼓或罷，或泣或歌。",
      "六四：月幾望，馬匹亡，無咎。",
      "九五：有孚攣如，無咎。",
      "上九：翰音登于天，貞凶。"
    ]
  },
  {
    "number": 62, "name": "小過", "upper": 1, "lower": 4,
    "judgment": "亨，利貞，可小事，不可大事。飛鳥遺之音，不宜上，宜下，大吉。",
    "lines": [
      "初六：飛鳥以凶。",
      "六二：過其祖，遇其妣；不及其君，遇其臣，無咎。",
      "九三：弗過防之，從或戕之，凶。",
      "九四：無咎，弗過遇之，往厲必戒，勿用永貞。",
      "六五：密雲不雨，自我西郊，公弋取彼在穴。",
      "上六：弗遇過之，飛鳥離之，凶，是謂災眚。"
    ]
  },
  {
    "number": 63, "name": "既濟", "upper": 2, "lower": 5,
    "judgment": "亨小，利貞，初吉終亂。",
    "lines": [
      "初九：曳其輪，濡其尾，無咎。",
      "六二：婦喪其茀，勿逐，七日得。",
      "九三：高宗伐鬼方，三年克之，小人勿用。",
      "六四：繻有衣袽，終日戒。",
      "九五：東鄰殺牛，不如西鄰之禴祭，實受其福。",
      "上六：濡其首，厲。"
    ]
  },
  {
    "number": 64, "name": "未濟", "upper": 5, "lower": 2,
    "judgment": "亨，小狐汔濟，濡其尾，無攸利。",
    "lines": [
      "初六：濡其尾，吝。",
      "九二：曳其輪，貞吉。",
      "六三：未濟，征凶，利涉大川。",
      "九四：貞吉，悔亡，震用伐鬼方，三年有賞于大國。",
      "六五：貞吉，無悔，君子之光，有孚，吉。",
      "上九：有孚于飲酒，無咎，濡其首，有孚失是。"
    ]
  }
]
```

- [ ] **Step 4：驗證 JSON 格式正確**

```bash
node -e "const d=require('./core-spec/hexagrams.json'); console.log('卦數：'+d.length); const pairs=new Set(d.map(h=>h.upper+','+h.lower)); console.log('唯一上下卦組合：'+pairs.size);"
```

預期輸出：
```
卦數：64
唯一上下卦組合：64
```

- [ ] **Step 5：Commit**

```bash
git add core-spec/
git commit -m "feat: 新增 core-spec 目錄與完整 64 卦資料"
```

---

## Task 2：Go 專案初始化 + types.go

**Files:**
- Create: `tui-go/go.mod`
- Create: `tui-go/internal/model/types.go`

- [ ] **Step 1：建立目錄**

```bash
mkdir -p tui-go/internal/model tui-go/data tui-go/cmd/yaonomad-tui
```

- [ ] **Step 2：建立 go.mod**

```bash
cd tui-go && go mod init github.com/LinGaOscar/yaonomad && cd ..
```

- [ ] **Step 3：寫入 tui-go/internal/model/types.go**

```go
package model

// Line 代表一次三硬幣擲出的爻結果
type Line struct {
	Value    int  // 三枚硬幣加總：6/7/8/9
	Yin      bool // true = 陰爻（偶數）
	Changing bool // true = 老陰(6) 或 老陽(9)，為變爻
}

// Hexagram 代表 64 卦之一
type Hexagram struct {
	Number   int
	Name     string
	Upper    int       // 上卦三位碼 0-7
	Lower    int       // 下卦三位碼 0-7
	Judgment string
	Lines    [6]string // index 0 = 初爻
}

// CastResult 代表一次完整卜卦結果
type CastResult struct {
	Question string
	Lines    [6]Line
	Present  Hexagram
	Future   *Hexagram // 無變爻時為 nil
}
```

- [ ] **Step 4：確認編譯無誤**

```bash
cd tui-go && go build ./... && cd ..
```

預期：無輸出（成功）。

- [ ] **Step 5：Commit**

```bash
git add tui-go/
git commit -m "feat: Go 專案初始化，定義 Line/Hexagram/CastResult 型別"
```

---

## Task 3：Go hexagrams.go（TDD）

**Files:**
- Create: `tui-go/internal/model/hexagrams_test.go`
- Create: `tui-go/internal/model/hexagrams.go`
- Copy: `tui-go/data/hexagrams.json`

- [ ] **Step 1：複製資料檔**

```bash
cp core-spec/hexagrams.json tui-go/data/hexagrams.json
```

- [ ] **Step 2：寫入 tui-go/internal/model/hexagrams_test.go**

```go
package model

import (
	"testing"
)

func loadTestHexagrams(t *testing.T) []Hexagram {
	t.Helper()
	hexagrams, err := LoadHexagrams("../../data/hexagrams.json")
	if err != nil {
		t.Fatalf("載入卦象失敗: %v", err)
	}
	return hexagrams
}

func TestLoadHexagrams_count(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	if len(hexagrams) != 64 {
		t.Errorf("卦象數量應為 64，實得 %d", len(hexagrams))
	}
}

func TestLoadHexagrams_uniquePairs(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	seen := make(map[[2]int]bool)
	for _, h := range hexagrams {
		key := [2]int{h.Upper, h.Lower}
		if seen[key] {
			t.Errorf("重複的上下卦組合: upper=%d lower=%d", h.Upper, h.Lower)
		}
		seen[key] = true
	}
}

func TestLookup_found(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	h, err := Lookup(hexagrams, 7, 7)
	if err != nil {
		t.Fatalf("查乾卦失敗: %v", err)
	}
	if h.Name != "乾" || h.Number != 1 {
		t.Errorf("乾卦資料錯誤: name=%s number=%d", h.Name, h.Number)
	}
}

func TestLookup_notFound(t *testing.T) {
	hexagrams := loadTestHexagrams(t)
	_, err := Lookup(hexagrams, 8, 8)
	if err == nil {
		t.Error("查詢無效碼（8）應回傳錯誤")
	}
}
```

- [ ] **Step 3：確認測試失敗（LoadHexagrams 未定義）**

```bash
cd tui-go && go test ./internal/model/ && cd ..
```

預期：編譯錯誤 `undefined: LoadHexagrams`。

- [ ] **Step 4：寫入 tui-go/internal/model/hexagrams.go**

```go
package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type hexagramJSON struct {
	Number   int       `json:"number"`
	Name     string    `json:"name"`
	Upper    int       `json:"upper"`
	Lower    int       `json:"lower"`
	Judgment string    `json:"judgment"`
	Lines    [6]string `json:"lines"`
}

// LoadHexagrams 從 JSON 檔案載入 64 卦資料
func LoadHexagrams(path string) ([]Hexagram, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀取卦象資料失敗: %w", err)
	}
	var entries []hexagramJSON
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("解析卦象資料失敗: %w", err)
	}
	hexagrams := make([]Hexagram, len(entries))
	for i, e := range entries {
		hexagrams[i] = Hexagram{
			Number:   e.Number,
			Name:     e.Name,
			Upper:    e.Upper,
			Lower:    e.Lower,
			Judgment: e.Judgment,
			Lines:    e.Lines,
		}
	}
	return hexagrams, nil
}

// Lookup 在卦象列表中查找指定上下卦的卦象
func Lookup(hexagrams []Hexagram, upper, lower int) (Hexagram, error) {
	for _, h := range hexagrams {
		if h.Upper == upper && h.Lower == lower {
			return h, nil
		}
	}
	return Hexagram{}, fmt.Errorf("找不到卦象: upper=%d lower=%d", upper, lower)
}
```

- [ ] **Step 5：確認測試通過**

```bash
cd tui-go && go test ./internal/model/ -v -run TestLoad && cd ..
```

預期：
```
--- PASS: TestLoadHexagrams_count (0.00s)
--- PASS: TestLoadHexagrams_uniquePairs (0.00s)
--- PASS: TestLookup_found (0.00s)
--- PASS: TestLookup_notFound (0.00s)
PASS
```

- [ ] **Step 6：Commit**

```bash
git add tui-go/
git commit -m "feat: Go hexagrams.go - 載入與查卦功能，含測試"
```

---

## Task 4：Go casting.go（TDD）

**Files:**
- Create: `tui-go/internal/model/casting_test.go`
- Create: `tui-go/internal/model/casting.go`

- [ ] **Step 1：寫入 tui-go/internal/model/casting_test.go**

```go
package model

import (
	"math/rand"
	"testing"
)

// fixedSource 提供固定輸出的隨機源，用於可重現測試
type fixedSource struct {
	seq []int64
	idx int
}

func (fs *fixedSource) Int63() int64 {
	v := fs.seq[fs.idx%len(fs.seq)]
	fs.idx++
	return v
}
func (fs *fixedSource) Seed(int64) {}

// newFixed 建立固定輸出的 *rand.Rand
// val=1 → Intn(2)=1 → coin()=3（背面）
// val=0 → Intn(2)=0 → coin()=2（正面）
func newFixed(val int64) *rand.Rand {
	return rand.New(&fixedSource{seq: []int64{val}})
}

func TestCastHexagram_allTails(t *testing.T) {
	// 全背面(3)：sum=9 老陽 → 本卦乾(#1)，變卦坤(#2)
	hexagrams := loadTestHexagrams(t)
	result, err := CastHexagram("測試", hexagrams, newFixed(1))
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		if l.Value != 9 {
			t.Errorf("爻 %d 應為 9，實得 %d", i, l.Value)
		}
		if l.Yin {
			t.Errorf("爻 %d 應為陽爻（Yin=false）", i)
		}
		if !l.Changing {
			t.Errorf("爻 %d 應為變爻（Changing=true）", i)
		}
	}
	if result.Present.Number != 1 {
		t.Errorf("本卦應為乾(#1)，實得 #%d %s", result.Present.Number, result.Present.Name)
	}
	if result.Future == nil {
		t.Fatal("六爻全為老陽，應有變卦")
	}
	if result.Future.Number != 2 {
		t.Errorf("變卦應為坤(#2)，實得 #%d %s", result.Future.Number, result.Future.Name)
	}
}

func TestCastHexagram_allHeads(t *testing.T) {
	// 全正面(2)：sum=6 老陰 → 本卦坤(#2)，變卦乾(#1)
	hexagrams := loadTestHexagrams(t)
	result, err := CastHexagram("測試", hexagrams, newFixed(0))
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	if result.Present.Number != 2 {
		t.Errorf("本卦應為坤(#2)，實得 #%d", result.Present.Number)
	}
	if result.Future == nil || result.Future.Number != 1 {
		t.Error("變卦應為乾(#1)")
	}
}

func TestCastHexagram_noChanging(t *testing.T) {
	// seq [0,0,1] → coins 2,2,3 → sum=7 少陽，重複六次
	hexagrams := loadTestHexagrams(t)
	r := rand.New(&fixedSource{seq: []int64{0, 0, 1}})
	result, err := CastHexagram("測試", hexagrams, r)
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		if l.Value != 7 {
			t.Errorf("爻 %d 應為 7，實得 %d", i, l.Value)
		}
		if l.Changing {
			t.Errorf("爻 %d 不應為變爻", i)
		}
	}
	if result.Future != nil {
		t.Error("無變爻時 Future 應為 nil")
	}
}

func TestCastHexagram_lineFields(t *testing.T) {
	// 驗證 Yin 與 Changing 欄位與 Value 一致
	hexagrams := loadTestHexagrams(t)
	r := rand.New(rand.NewSource(42))
	result, err := CastHexagram("測試", hexagrams, r)
	if err != nil {
		t.Fatalf("卜卦失敗: %v", err)
	}
	for i, l := range result.Lines {
		wantYin := l.Value%2 == 0
		if l.Yin != wantYin {
			t.Errorf("爻 %d Yin=%v 應為 %v（value=%d）", i, l.Yin, wantYin, l.Value)
		}
		wantChanging := l.Value == 6 || l.Value == 9
		if l.Changing != wantChanging {
			t.Errorf("爻 %d Changing=%v 應為 %v（value=%d）", i, l.Changing, wantChanging, l.Value)
		}
	}
}
```

- [ ] **Step 2：確認測試失敗（CastHexagram 未定義）**

```bash
cd tui-go && go test ./internal/model/ -run TestCast && cd ..
```

預期：編譯錯誤 `undefined: CastHexagram`。

- [ ] **Step 3：寫入 tui-go/internal/model/casting.go**

```go
package model

import "math/rand"

// CastHexagram 執行一次三硬幣卜卦
// r 為隨機源，測試時注入固定種子以確保可重現
func CastHexagram(question string, hexagrams []Hexagram, r *rand.Rand) (CastResult, error) {
	var lines [6]Line
	for i := range lines {
		sum := coin(r) + coin(r) + coin(r)
		lines[i] = Line{
			Value:    sum,
			Yin:      sum%2 == 0,
			Changing: sum == 6 || sum == 9,
		}
	}

	upper, lower := trigramCodes(lines)
	present, err := Lookup(hexagrams, upper, lower)
	if err != nil {
		return CastResult{}, err
	}

	future, err := computeFuture(hexagrams, lines)
	if err != nil {
		return CastResult{}, err
	}

	return CastResult{
		Question: question,
		Lines:    lines,
		Present:  present,
		Future:   future,
	}, nil
}

// coin 模擬一枚硬幣：正面(字面)=2，背面=3
func coin(r *rand.Rand) int {
	if r.Intn(2) == 0 {
		return 2
	}
	return 3
}

// trigramCodes 計算六爻的上下卦三位碼
// bit0 = 底爻，yang=1，yin=0
func trigramCodes(lines [6]Line) (upper, lower int) {
	yang := func(l Line) int {
		if l.Yin {
			return 0
		}
		return 1
	}
	lower = yang(lines[0]) | yang(lines[1])<<1 | yang(lines[2])<<2
	upper = yang(lines[3]) | yang(lines[4])<<1 | yang(lines[5])<<2
	return upper, lower
}

// computeFuture 計算變卦；若無變爻則回傳 nil
func computeFuture(hexagrams []Hexagram, lines [6]Line) (*Hexagram, error) {
	hasChanging := false
	for _, l := range lines {
		if l.Changing {
			hasChanging = true
			break
		}
	}
	if !hasChanging {
		return nil, nil
	}

	var futureLines [6]Line
	for i, l := range lines {
		yin := l.Yin
		if l.Changing {
			yin = !l.Yin
		}
		futureLines[i] = Line{Value: l.Value, Yin: yin, Changing: l.Changing}
	}

	upper, lower := trigramCodes(futureLines)
	h, err := Lookup(hexagrams, upper, lower)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
```

- [ ] **Step 4：確認所有測試通過**

```bash
cd tui-go && go test ./internal/model/ -v && cd ..
```

預期：所有測試 PASS，無 FAIL。

- [ ] **Step 5：Commit**

```bash
git add tui-go/internal/model/casting.go tui-go/internal/model/casting_test.go
git commit -m "feat: Go casting.go - 三硬幣卜卦核心，含測試"
```

---

## Task 5：Web core 目錄 + models.js

**Files:**
- Create: `web/core/models.js`
- Create: `web/data/hexagrams.json`

- [ ] **Step 1：建立目錄**

```bash
mkdir -p web/core web/data
```

- [ ] **Step 2：複製資料檔**

```bash
cp core-spec/hexagrams.json web/data/hexagrams.json
```

- [ ] **Step 3：建立 web/package.json（啟用 ES Modules）**

```json
{
  "type": "module"
}
```

寫入路徑：`web/package.json`。Node.js 執行 `node --test` 時需要此檔案才能識別 `import`/`export` 語法。

- [ ] **Step 4：寫入 web/core/models.js**

```js

/**
 * @typedef {Object} Line
 * @property {number} value - 爻值，固定為 6、7、8 或 9
 * @property {boolean} yin - true = 陰爻（value 為偶數）
 * @property {boolean} changing - true = 老陰(6) 或 老陽(9)
 */

/**
 * @typedef {Object} Hexagram
 * @property {number} number - 文王序列 1–64
 * @property {string} name - 卦名
 * @property {number} upper - 上卦三位碼（0–7）
 * @property {number} lower - 下卦三位碼（0–7）
 * @property {string} judgment - 卦辭
 * @property {string[]} lines - 爻辭，index 0 = 初爻，共 6 條
 */

/**
 * @typedef {Object} CastResult
 * @property {string} question - 問題
 * @property {Line[]} lines - 六爻（index 0 = 初爻）
 * @property {Hexagram} present - 本卦
 * @property {Hexagram|null} future - 變卦；無變爻時為 null
 */

// 此檔案僅提供 JSDoc 型別定義，無執行時期程式碼
```

- [ ] **Step 5：Commit**

```bash
git add web/
git commit -m "feat: Web 核心目錄初始化，JSDoc 型別定義"
```

---

## Task 6：Web hexagrams.js（TDD）

**Files:**
- Create: `web/core/hexagrams.test.js`
- Create: `web/core/hexagrams.js`

- [ ] **Step 1：寫入 web/core/hexagrams.test.js**

```js
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { lookup } from './hexagrams.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
const hexagrams = JSON.parse(readFileSync(join(__dirname, '../data/hexagrams.json'), 'utf-8'))

test('hexagrams.json 應有 64 筆', () => {
  assert.strictEqual(hexagrams.length, 64)
})

test('hexagrams.json 的 upper/lower 組合應唯一', () => {
  const seen = new Set()
  for (const h of hexagrams) {
    const key = `${h.upper},${h.lower}`
    assert.ok(!seen.has(key), `重複的組合: ${key}`)
    seen.add(key)
  }
})

test('lookup 找到乾卦 (upper=7, lower=7)', () => {
  const h = lookup(hexagrams, 7, 7)
  assert.strictEqual(h.number, 1)
  assert.strictEqual(h.name, '乾')
})

test('lookup 找到坤卦 (upper=0, lower=0)', () => {
  const h = lookup(hexagrams, 0, 0)
  assert.strictEqual(h.number, 2)
  assert.strictEqual(h.name, '坤')
})

test('lookup 查無效碼應拋出錯誤', () => {
  assert.throws(() => lookup(hexagrams, 8, 8), /找不到卦象/)
})
```

- [ ] **Step 2：確認測試失敗（hexagrams.js 不存在）**

```bash
node --test web/core/hexagrams.test.js
```

預期：`Error: Cannot find module './hexagrams.js'`

- [ ] **Step 3：寫入 web/core/hexagrams.js**

```js
'use strict'

/**
 * 在卦象列表中以上下卦三位碼查找對應卦象
 * @param {import('./models.js').Hexagram[]} hexagrams
 * @param {number} upper - 上卦三位碼（0–7）
 * @param {number} lower - 下卦三位碼（0–7）
 * @returns {import('./models.js').Hexagram}
 * @throws {Error} 若找不到對應卦象
 */
export function lookup(hexagrams, upper, lower) {
  const found = hexagrams.find(h => h.upper === upper && h.lower === lower)
  if (!found) {
    throw new Error(`找不到卦象: upper=${upper} lower=${lower}`)
  }
  return found
}
```

- [ ] **Step 4：確認測試通過**

```bash
node --test web/core/hexagrams.test.js
```

預期：
```
✓ hexagrams.json 應有 64 筆
✓ hexagrams.json 的 upper/lower 組合應唯一
✓ lookup 找到乾卦 (upper=7, lower=7)
✓ lookup 找到坤卦 (upper=0, lower=0)
✓ lookup 查無效碼應拋出錯誤
```

- [ ] **Step 5：Commit**

```bash
git add web/core/hexagrams.js web/core/hexagrams.test.js
git commit -m "feat: Web hexagrams.js - lookup 函式，含測試"
```

---

## Task 7：Web casting.js（TDD）

**Files:**
- Create: `web/core/casting.test.js`
- Create: `web/core/casting.js`

- [ ] **Step 1：寫入 web/core/casting.test.js**

```js
import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { castHexagram } from './casting.js'

const __dirname = dirname(fileURLToPath(import.meta.url))
const hexagrams = JSON.parse(readFileSync(join(__dirname, '../data/hexagrams.json'), 'utf-8'))

test('全背面(3) → 本卦乾(#1)，變卦坤(#2)', () => {
  const result = castHexagram(hexagrams, '測試', () => 3)
  assert.strictEqual(result.present.number, 1)
  assert.strictEqual(result.present.name, '乾')
  assert.ok(result.future !== null, '應有變卦')
  assert.strictEqual(result.future.number, 2)
  assert.strictEqual(result.future.name, '坤')
})

test('全正面(2) → 本卦坤(#2)，變卦乾(#1)', () => {
  const result = castHexagram(hexagrams, '測試', () => 2)
  assert.strictEqual(result.present.number, 2)
  assert.strictEqual(result.present.name, '坤')
  assert.ok(result.future !== null, '應有變卦')
  assert.strictEqual(result.future.number, 1)
  assert.strictEqual(result.future.name, '乾')
})

test('少陽序列(2,2,3=7) → 無變爻 → future 為 null', () => {
  const coins = [2, 2, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++ % 3])
  for (const l of result.lines) {
    assert.strictEqual(l.value, 7)
    assert.strictEqual(l.yin, false)
    assert.strictEqual(l.changing, false)
  }
  assert.strictEqual(result.future, null)
})

test('少陰序列(2,3,3=8) → 無變爻 → future 為 null', () => {
  const coins = [2, 3, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++ % 3])
  for (const l of result.lines) {
    assert.strictEqual(l.value, 8)
    assert.strictEqual(l.yin, true)
    assert.strictEqual(l.changing, false)
  }
  assert.strictEqual(result.future, null)
})

test('question 欄位正確帶入', () => {
  const result = castHexagram(hexagrams, '今日運勢如何？', () => 3)
  assert.strictEqual(result.question, '今日運勢如何？')
})

test('爻值均在合法範圍 6–9', () => {
  const coins = [3, 3, 3, 2, 2, 2, 3, 2, 3, 2, 3, 2, 3, 3, 2, 2, 2, 3]
  let i = 0
  const result = castHexagram(hexagrams, '測試', () => coins[i++])
  for (const l of result.lines) {
    assert.ok([6, 7, 8, 9].includes(l.value), `爻值 ${l.value} 不合法`)
  }
})
```

- [ ] **Step 2：確認測試失敗（casting.js 不存在）**

```bash
node --test web/core/casting.test.js
```

預期：`Error: Cannot find module './casting.js'`

- [ ] **Step 3：寫入 web/core/casting.js**

```js
'use strict'

import { lookup } from './hexagrams.js'

/**
 * 執行一次三硬幣卜卦
 * @param {import('./models.js').Hexagram[]} hexagrams - 預先載入的 64 卦資料
 * @param {string} question - 使用者問題
 * @param {() => number} [coinFn] - 回傳 2 或 3 的函式；預設以 Math.random 模擬
 * @returns {import('./models.js').CastResult}
 */
export function castHexagram(hexagrams, question, coinFn = () => (Math.random() < 0.5 ? 2 : 3)) {
  const lines = []
  for (let i = 0; i < 6; i++) {
    const sum = coinFn() + coinFn() + coinFn()
    lines.push({
      value: sum,
      yin: sum % 2 === 0,
      changing: sum === 6 || sum === 9,
    })
  }

  const lower = trigramCode(lines.slice(0, 3))
  const upper = trigramCode(lines.slice(3, 6))
  const present = lookup(hexagrams, upper, lower)

  const future = lines.some(l => l.changing)
    ? computeFuture(hexagrams, lines)
    : null

  return { question, lines, present, future }
}

/**
 * 計算三條爻的 trigram 三位碼（bit0 = 底爻，yang=1，yin=0）
 * @param {import('./models.js').Line[]} threeLines
 * @returns {number}
 */
function trigramCode(threeLines) {
  return threeLines.reduce((code, line, i) => {
    const yang = line.yin ? 0 : 1
    return code | (yang << i)
  }, 0)
}

/**
 * 計算變卦（將變爻的陰陽反轉後查卦）
 * @param {import('./models.js').Hexagram[]} hexagrams
 * @param {import('./models.js').Line[]} lines
 * @returns {import('./models.js').Hexagram}
 */
function computeFuture(hexagrams, lines) {
  const futureLines = lines.map(l => ({
    value: l.value,
    yin: l.changing ? !l.yin : l.yin,
    changing: l.changing,
  }))
  const lower = trigramCode(futureLines.slice(0, 3))
  const upper = trigramCode(futureLines.slice(3, 6))
  return lookup(hexagrams, upper, lower)
}
```

- [ ] **Step 4：確認所有測試通過**

```bash
node --test web/core/casting.test.js
```

預期：全部 6 個測試 PASS。

- [ ] **Step 5：執行所有 Web 測試確認無迴歸**

```bash
node --test web/core/hexagrams.test.js web/core/casting.test.js
```

預期：11 個測試全部 PASS。

- [ ] **Step 6：執行所有 Go 測試確認無迴歸**

```bash
cd tui-go && go test ./... && cd ..
```

預期：PASS。

- [ ] **Step 7：Commit**

```bash
git add web/core/casting.js web/core/casting.test.js
git commit -m "feat: Web casting.js - 三硬幣卜卦核心，含測試"
```

---

## 完成確認

所有任務完成後執行最終驗證：

```bash
# Go 全部測試
cd tui-go && go test ./... -v && cd ..

# Web 全部測試
node --test web/core/hexagrams.test.js web/core/casting.test.js
```

預期：Go 8 個測試 PASS，Web 11 個測試 PASS，無任何 FAIL。
