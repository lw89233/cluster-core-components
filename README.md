# Cluster Core Components

Zestaw wyizolowanych modułów infrastrukturalnych w języku Go, służących do zarządzania stanem, harmonogramowania i rekoncyliacji w środowiskach rozproszonych.

## Architektura Systemu

Projekt składa się z następujących, niezależnych pakietów:

* **pkg/storage** - Moduł zapewniający atomowy zapis i odczyt zrzutów pamięci (JSON) z wykorzystaniem muteksów. Gwarantuje monotoniczne wersjonowanie, zapobiegając nadpisywaniu nowszego stanu starszymi plikami (Crash Recovery).
* **pkg/scheduler** - Zaawansowany planista przydzielający zadania do węzłów roboczych. Wykorzystuje system punktacji (Scoring), sprawdza pojemności abstrakcyjne (CPU/RAM) oraz wdraża reguły Anti-Affinity dla zadań bezstanowych.
* **pkg/reconciler** - Silnik różnicujący (Diffing Engine). Oblicza wymagane akcje naprawcze (CREATE/DELETE) na podstawie porównania stanu pożądanego (Desired State) z aktualnym (Actual State).
* **pkg/liveness** - Współbieżny monitor aktywności węzłów. Automatycznie degraduje statusy (ACTIVE -> STALE -> DEAD) na podstawie zdefiniowanych progów czasowych braku aktywności.
* **pkg/dispatcher** - Moduł asynchronicznej dystrybucji zdarzeń, zrealizowany w oparciu o wzorzec Obserwatora i bezblokujące kanały komunikacyjne (channels) języka Go.

## Profilowanie i Wydajność

Moduły zostały zoptymalizowane pod kątem niskiego zużycia zasobów i dużej współbieżności. Rdzeń harmonogramujący potrafi zaalokować 5000 zadań na 100 węzłach w czasie poniżej 80 milisekund (przy architekturze x86_64).