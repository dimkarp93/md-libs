# md-libs

Общее ядро для конвертеров markdown: парсинг, фильтрация по заголовкам, страничный пайплайн и контракт рендеринга. Используется в [md-docx](https://github.com/dimkarp93/md-docx) и [md-pdf](https://github.com/dimkarp93/md-pdf).

Библиотека не тянет внешних зависимостей и не содержит конкретных бэкендов вывода — их реализуют потребители.

## Установка

```sh
go get github.com/dimkarp93/md-libs
```

## Пакеты

### `mdlib` — пайплайн

Полный путь от исходного markdown до отфильтрованного списка блоков.

```go
blocks, err := mdlib.Prepare(src, mdlib.Options{
    Pages:            "1,3-5", // пусто — все страницы, кроме нулевой
    Heads:            "h2:result,h3:resume",
    HideMatchedHeads: true,
})
```

### `markdown` — парсинг

```go
type Block struct {
    Kind  BlockKind // Paragraph | Heading | Code | PageBreak
    Level int
    Text  string
    Lines []string
}

func Parse(src string) []Block
func ParseInline(text string, allowUnderscore bool) []Run
func PlainText(text string) string
func SplitPages(src string) []string
func ParsePageRanges(spec string) ([]int, error)
```

Документ делится на страницы по строкам `---` (внутри код-фенса разделитель игнорируется). Если файл начинается с закрытого блока `---...---`, его содержимое становится **нулевой страницей** (front matter) и по умолчанию в вывод не попадает.

### `filter` — отбор по заголовкам

```go
filters, err := filter.ParseHeads("h2:result,resume")
blocks = filter.ByHeads(blocks, filters, hideMatched)
```

Спецификация `h<N>:<имя>` привязывает фильтр к уровню заголовка; голое имя матчит любой уровень. Сравнение регистронезависимое и по «чистому» тексту, поэтому `**Result**` матчится фильтром `h2:result`. В вывод попадает заголовок и всё его содержимое до следующего заголовка того же или более высокого уровня. `hideMatched` убирает сам заголовок, оставляя содержимое.

### `render` — контракт рендеринга

```go
type Renderer interface {
    Heading(level int, text string)
    Paragraph(text string)
    Code(lines []string)
    PageBreak()
}

func Blocks(r Renderer, blocks []markdown.Block)
func HeadingSizePt(level int) float64
```

Бэкенд реализует четыре метода, `render.Blocks` раскладывает по ним список блоков. `HeadingSizePt` — единая шкала размеров заголовков (18/16/14/13/12/11 pt), чтобы docx и pdf не расходились.

## Версионирование

Версия хранится в `versions.txt`. Push в `main` с новой версией запускает workflow, который проверяет тесты и создаёт релиз с тегом `vX.Y.Z` — это ровно тот формат, который нужен Go-модулю.

```sh
make bump-patch   # или bump-minor / bump-major
```

Пока API не устоялся, серия остаётся `v0.x`: в нулевой мажорной версии ломающие изменения допустимы без суффикса `/v2` в module path.

## Разработка

```sh
make build                        # компиляция всех пакетов
make test                         # тесты
make test-v                       # то же, с именами тестов
make test-run T=TestParseInline   # один тест или маска
make cover                        # покрытие + итоговый процент
make vet
make check                        # vet + test
```

`make test-all` дополнительно прогоняет тесты в `../md-docx` и `../md-pdf`, если они склонированы рядом, — одна команда проверяет и библиотеку, и обоих потребителей.

Что где тестируется:

- **здесь** — парсинг (`markdown`), фильтрация (`filter`), пайплайн (`mdlib`) и контракт рендеринга (`render`): то, что общее для всех потребителей;
- **в md-docx и md-pdf** — конкретные рендереры (XML для Word, вывод fpdf) и сквозные тесты CLI: то, что у каждого своё.

Потребители, которым нужно править библиотеку и CLI одновременно, подключают локальную копию через `go.work` — см. `make configure` в md-docx и md-pdf.

## Пакет для ручного переноса

`make pack` собирает библиотеку в виде **файлового модуль-прокси** — того же формата, что отдаёт `proxy.golang.org`:

```sh
make pack
# dist/md-libs-0.1.0-proxy.tar.gz
```

Внутри — `github.com/dimkarp93/md-libs/@v/` с файлами `v0.1.0.zip`, `v0.1.0.mod`, `v0.1.0.info` и `list`. Это позволяет собирать зависимый проект против обычного `require github.com/dimkarp93/md-libs v0.1.0` — **без `replace` и без правки `go.mod`**:

```sh
tar -xzf md-libs-0.1.0-proxy.tar.gz -C /opt
cd ../md-docx
GOFLAGS=-mod=mod GOPROXY=file:///opt/proxy GOSUMDB=off go build ./...
```

Так собирается проект, у которого нет других внешних зависимостей (md-docx) или все они уже в кэше модулей. Если остальные зависимости нужно тянуть из сети (md-pdf зависит от `fpdf`), добавьте сеть в цепочку прокси и выключайте sumdb **точечно**:

```sh
GOFLAGS=-mod=mod \
GOPROXY=file:///opt/proxy,https://proxy.golang.org,direct \
GONOSUMDB='github.com/dimkarp93/*' \
go build ./...
```

- `GOSUMDB=off` / `GONOSUMDB` обязателен, пока модуль не опубликован: иначе Go пойдёт в `sum.golang.org` и получит 404. Глобальный `GOSUMDB=off` при этом ломает проверку остальных модулей, скачиваемых из сети, — поэтому в смешанном сценарии нужен именно `GONOSUMDB`.
- **Не используйте `GOPRIVATE`**: он попутно выставляет `GONOPROXY`, и Go пойдёт за md-libs напрямую в GitHub мимо файлового прокси.
- `GOFLAGS=-mod=mod` нужен, чтобы Go дописал `go.sum`; после этого обычная сборка работает и без него.
- Если в проекте включён workspace (`go.work`), добавьте `GOWORK=off` — иначе он перекроет прокси.

Отличие от `make configure` в CLI: `configure` подсовывает библиотеку с диска через `go.work` и удобен, когда обе части правятся одновременно; `pack` даёт самодостаточный артефакт с конкретной версией — для машины без доступа к репозиторию.

Пакет собирается из рабочего дерева как есть, поэтому запускайте `make pack` на чистом дереве: содержимое zip определяет контрольную сумму в `go.sum`, и она должна совпасть с той, что позже посчитает GitHub для тега `vX.Y.Z`.

Один внешний ресурс всё же нужен: `go.mod` объявляет `go 1.26.1`, и Go докачивает такой toolchain как модуль. На машине без сети он должен быть уже установлен (либо соберите с `GOTOOLCHAIN=local` достаточно свежим Go).
