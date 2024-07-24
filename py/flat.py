def _as_flat_list(d, parent_key=''):
    items = []
    for k, v in d.items():
        new_key = parent_key + '.' + k if parent_key else k
        if isinstance(v, dict):
            items.extend(_as_flat_list(v, new_key))
        else:
            items.append((new_key, v))
    return items


def flat_dict(d: dict) -> dict:
    """Flattens the dictionary key using "." to join nested dictionaries.

    The function is useful for MongoDB updates to ensure it doesn't rewrite nested objects, but
    rather updates their fields.

    >>> flat_dict({})
    {}
    >>> flat_dict({'a': 'b'})
    {'a': 'b'}
    >>> flat_dict({'a': 'b', 'c': 'd'})
    {'a': 'b', 'c': 'd'}
    >>> flat_dict({'a': {'b': 'c'}})
    {'a.b': 'c'}
    >>> flat_dict({'a': 'b', 'c': {'d': 'e', 'f': {'g': 'h', 'i': 'j'}, 'k': 'l'}})
    {'a': 'b', 'c.d': 'e', 'c.f.g': 'h', 'c.f.i': 'j', 'c.k': 'l'}
    """
    return dict(_as_flat_list(d))


if __name__ == "__main__":
    import doctest
    doctest.testmod()
