---
title: "Python 코딩테스트 기초 Tip 모음"
date: 2025-10-28
categories:
  - python
tags:
  - python
  - 코딩테스트
  - 알고리즘
  - 자료구조
---

Python 코딩테스트를 준비하면서 알아두면 유용한 문법과 기법들을 정리했습니다.

## 1. 문자열 뒤집기

`[::-1]` 슬라이싱으로 간단하게 뒤집을 수 있습니다.

```python
string = 'Hello, World!'
print(string[::-1])  # !dlroW ,olleH
```

## 2. 중복 제거하기

`set` 자료형은 중복을 허용하지 않고 순서가 없습니다.

```python
temp = [1, 1, 2, 2, 3, 4, 4, 5]
print(list(set(temp)))  # [1, 2, 3, 4, 5]
```

## 3. 한 줄에 여러 정수 입력받기

`map()` 함수와 `split()`을 조합합니다.

```python
temp_list = list(map(int, input().split()))
# 입력: 34 2 566 4 7 8 11
# 출력: [34, 2, 566, 4, 7, 8, 11]
```

## 4. 2차원 맵 생성

BFS, DFS 등에서 방문 여부를 체크할 때 사용합니다.

```python
visited = [[False for _ in range(m)] for _ in range(n)]
```

## 5. 연속 비교 연산

Python은 `if 0 < n < 10:` 같은 연속 비교를 허용합니다.

```python
# Java/C: if (0 < n && n < 10) { }
# Python:
if 0 < n < 10:
    pass
```

## 6. 두 변수의 값 바꾸기

```python
a, b = b, a
```

## 7. enumerate()로 인덱스와 값 동시에 가져오기

```python
temp = ['k', 'o', 'r', 'e', 'a']
for idx, value in enumerate(temp):
    print(idx, value)
```

## 8. deque 사용하기 (필수!)

`pop(0)`은 O(N)이지만 `deque.popleft()`는 **O(1)**입니다. BFS 구현 시 반드시 사용됩니다.

```python
from collections import deque

queue = deque([1, 2, 3, 4, 5])
print(queue.popleft())  # 1
```

## 9. zip()으로 여러 리스트 동시 순회

```python
temp1 = [1, 3, 5]
temp2 = [2, 4, 6]
for n1, n2 in zip(temp1, temp2):
    print(n1, n2)  # 1 2 / 3 4 / 5 6
```

## 10. 딕셔너리 정렬

```python
dic = {'apple': 3, 'banana': 1, 'pear': 5}

# value 기준 정렬
sorted(dic.items(), key=lambda x: x[1])

# key 기준 정렬
sorted(dic)
```

## 11. for-else / while-else 문

break 없이 루프가 끝나면 else 블록이 실행됩니다. flag 변수가 필요 없어 편합니다.

```python
for i in range(1, 10):
    if i == 11:
        break
else:
    print('break 안걸림!')  # 출력됨
```

## 12. bisect - 이진 탐색

정렬된 리스트에서 O(logN)으로 값의 위치를 찾습니다.

```python
import bisect

lst = [1, 3, 5, 6, 6, 8]
print(bisect.bisect_left(lst, 4))   # 2
print(bisect.bisect_left(lst, 6))   # 3
print(bisect.bisect_right(lst, 6))  # 5
```

## 13. 2차원 리스트에서 열 추출하기

```python
a = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
b = list(zip(*a))[0]  # (1, 4, 7)
```

## 14. sys.stdin.readline으로 입력 시간 단축

입력이 많을 때 필수적입니다.

```python
import sys
input = sys.stdin.readline

num = int(input())  # 정수 캐스팅 시 strip() 불필요
```

## 15. defaultdict - 키 자동 생성 딕셔너리

```python
from collections import defaultdict
dic = defaultdict(list)

dic[1].append('temp')  # KeyError 없이 자동으로 키 생성
```

## 16. 우선순위 큐 (heapq)

push/pop 모두 **O(logN)**으로 매우 효율적입니다. 그리디/정렬 골드 이상 문제에서 먼저 떠올려보세요.

```python
import heapq

pq = []
heapq.heappush(pq, (3, '작업1'))
heapq.heappush(pq, (1, '작업2'))

# 최소 힙 (기본)
priority, task = heapq.heappop(pq)

# 최대 힙: 우선순위에 -를 곱해서 넣기
heapq.heappush(pq, (-priority, task))
```

## 17. 2차원 리스트 깊은복사

`[:]`는 1차원만 깊은복사됩니다. 2차원 이상에서는 `deepcopy()`를 사용해야 합니다.

```python
import copy

original = [[1, 2, 3], [4, 5, 6]]

# 얕은 복사 - 내부 리스트는 참조 공유!
shallow = original[:]
shallow[0][0] = 999
print(original[0][0])  # 999 (원본도 변경됨!)

# 깊은 복사 - 완전히 독립적
deep = copy.deepcopy(original)
deep[0][0] = 111
print(original[0][0])  # 999 (원본 유지)
```
