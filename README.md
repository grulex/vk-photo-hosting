# Бесплатный хостинг фотографий VK
Для хакатонов или пет-проектов.

## Установка
```shell
go get github.com/grulex/vk-photo-hosting
```

## Начало
- Регистрируемся на vk.com
- Создаём новое сообщество [тут](https://vk.com/groups?w=groups_create). 
Рекомендуется в "тип группы" выбирать "частная", чтобы случайные посетители не видели ваши альбомы
- В "Альбомах" создаём новый альбом. Открыть раздел можно по ссылке вида `https://vk.com/albums-{id сообщества}`
- Копируем URL полученного альбома. В ссылке видим 
`{Id группы}-{Id Альбома}`
- Получаем токен пользователя [тут](https://oauth.vk.com/authorize?client_id=6287487&scope=262148&redirect_uri=https://oauth.vk.com/blank.html&display=page&response_type=token&revoke=1). После предоставления доступа копируем из URL значение`access_token`

Таким образом, мы получили всё необходимое для работы с библиотекой: userToken, groupId, albumId

## Использование
### Создание объекта хостинга
```go
import "github.com/grulex/vk-photo-hosting"
//...
hosting := vkphotohosting.NewHosting(userToken, groupId, time.Minute)
```
### Загрузка файла в альбом
```go
photoPath := "/tmp/myphoto.jpg"
id, variants, err := hosting.UploadByFile(context.Background(), albumId, photoPath)
```
- id — идентификатор фотографии в альбоме
- variants — коллекция с разными размерами фотографии. Есть методы для получения максимального и минимального размера. В каждом объекте есть свой URL и размеры.
### Все методы
```go
// Загрузка из Reader
UploadByReader(ctx context.Context, albumId uint64, image io.Reader)

// Загрузка по пути к файлу
UploadByFile(ctx context.Context, albumId uint64, filePath string) (id uint64, variants internal.Variants, err error)

// Загрузка из внешнего URL в свой альбом
UploadByUrl(
    ctx context.Context,
    albumId uint64,
    photoUrl string,
    downloadTimeout time.Duration,
) (
    id uint64,
    variants internal.Variants,
    err error,
)
```
## Лимиты
- Количество запросов в секунду. Можно повторно выполнить через 1 сек.
- 5000 фотографий в одном альбоме. Можно создать новый альбом в группе
