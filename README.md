# TechFinder
TechFinder, verilen bir domain listesindeki web sitelerinin kullandıkları teknolojileri tespit etmek amacıyla geliştirilmiş bir analiz aracıdır.

### Nasıl Çalışır?
TechFinder, hedef web sitesinden elde ettiği çeşitli verileri analiz ederek kullanılan teknolojileri çıkarmaya çalışır:

- HTTP Response Header bilgileri üzerinden sunucu, framework veya diğer teknolojiler belirlenmeye çalışılır.

- HTML Body içeriğinden:

    - Yorum satırları taranarak bileşen adı ve sürüm numarası tespit edilir.

    - CSS ve JavaScript dosyaları gibi referanslar üzerinden regex ile teknoloji analizi yapılır.

- Cookie değerleri üzerinden teknolojilere dair ipuçları toplanır.

- Mevcut Özellikler
    - Header, body ve cookie analizine dayalı teknoloji tespiti.

    - Versiyon bilgisi ile birlikte bileşen tanımlama.

    - Yorum satırları üzerinden teknoloji tahmini.

- Yapılacaklar
    - HTTP yönlendirmeleri (3xx) otomatik olarak takip edilmeli.

    - [Snyk](https://snyk.io/) haricinde farklı güvenlik ve teknoloji analiz servisleri entegre edilmeli.