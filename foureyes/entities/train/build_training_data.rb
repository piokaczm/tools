require 'json'

class Entity
    attr_reader :text, :span, :answer

    def initialize(text, answer, span)
        @text = text
        @answer = answer
        @span = span
    end

    def to_json
        {
            "text": text,
            "answer": answer,
            "span": {
                "start": span.start_idx,
                "end": span.end_idx,
                "label": span.label,
            }
        }.to_json
    end
end

class Span
    attr_reader :start_idx, :end_idx, :label

    def initialize(start_idx, end_idx, label)
        @start_idx = start_idx
        @end_idx = end_idx
        @label = label
    end
end

def process(file, app_name, entity_name)
    entities = read_file(file, app_name, entity_name)
    
    File.open("#{app_name}.json", "w+") do |f|
        entities.each { |element| f.puts(element.to_json) }
    end
end

def read_file(file, app_name, entity_name)
    entities = []

    File.readlines(file).each do |line|
        if line.downcase.include?(app_name)
            s, e = get_indexes(line.downcase, app_name)
            span = Span.new(s, e, entity_name)
            entities << Entity.new(line, "accept", span)
        else
            word = line.split(" ").sample
            s, e = get_indexes(line.downcase, word.downcase)
            span = Span.new(s, e, entity_name)
            entities << Entity.new(line, "reject", span)
        end
    end

    entities
end

def get_indexes(line, word)
    idx = line.index(word)
    [idx, idx+word.size]
end

def main
    unless ARGV[0] && ARGV[1] && ARGV[2]
        puts "not enough args provided - 0-path to file, 1-name of the app"
        exit(1)
    end

    file_path = ARGV[0]
    app_name = ARGV[1]
    entity_name = ARGV[2]

    process(file_path, app_name, entity_name)
end

main()